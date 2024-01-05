package repositorys

import (
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"poll-application-service/models"
	"service-shared/database"
	sharedhttp "service-shared/http"
	messagequeue "service-shared/message-queue"
	sharedconfig "service-shared/shared-config"
	sharedmodels "service-shared/shared-models"
	"sync"
)

/*
RabbitMQWorker defines a worker responsible for handling
messages from a rabbitMQ queue. It receives Deliveries via
inChan and will attempt to process them.
*/
type RabbitMQWorker struct {
	repository      database.Repository
	wg              *sync.WaitGroup
	inChan          <-chan amqp.Delivery
	cfg             sharedconfig.Config
	deliveryHandler messagequeue.DeliveryHandler
	httpClient      sharedhttp.Client
}

func NewRabbitMQWorker(
	repo database.Repository,
	wg *sync.WaitGroup,
	inChan <-chan amqp.Delivery,
	cfg sharedconfig.Config,
	handler messagequeue.DeliveryHandler,
	httpClient sharedhttp.Client) RabbitMQWorker {
	return RabbitMQWorker{
		repository:      repo,
		wg:              wg,
		inChan:          inChan,
		cfg:             cfg,
		deliveryHandler: handler,
		httpClient:      httpClient,
	}
}

//ProcessMessages receives deliveries from its inChan and delegates
//processing responsibilities to processMessage(delivery amqp.Delivery).
func (worker RabbitMQWorker) ProcessMessages() {
	defer worker.wg.Done()
	for {
		delivery := <-worker.inChan
		worker.processMessage(delivery)
	}
}

/*
processMessage will process a delivery message.

This will reach out to the jobs API of the bank. If the status of an application
is still pending, it will be re-queued to be consumed again later.

If it is finished, ie the status is complete or rejected, then a call will be made
to update the db with the latest status.

If an unrecoverable error occurs while processing a message, this function will simulate
the behaviour of passing a message to a dead letter queue. Note that the DLQ behaviour has
not yet been implemented due to time constraints.
*/
func (worker RabbitMQWorker) processMessage(delivery amqp.Delivery) {
	body := delivery.Body

	var message *sharedmodels.PollLoanMessage
	err := json.Unmarshal(body, &message)
	if messagequeue.CheckError(
		err,
		fmt.Sprintf("Could not unmarshal message %s to PollLoanMessage - bad data on queue?\n", delivery.Body),
		delivery,
		worker.deliveryHandler) {
		return
	}

	finished, err := worker.pollApplicationStatus(message.BankApplicationID, message.OurApplicationID)
	if messagequeue.CheckError(err, "", delivery, worker.deliveryHandler) {
		// Something went wrong polling the status. Bank API might be down for example
		return
	}

	if !finished {
		// The loan is still pending so requeue
		worker.deliveryHandler.Nack(false, true, delivery)
		return
	}

	worker.deliveryHandler.Ack(false, delivery)
}

func (worker RabbitMQWorker) pollApplicationStatus(bankAppID string, ourApplicationID string) (bool, error) {
	pollRequest := models.PollLoanRequest{ApplicationID: bankAppID}
	resp, err := worker.sendPollRequest(pollRequest)
	if err != nil {
		log.Printf("Failed to send request to bank API %s\n", err)
		return false, err
	}

	finished, err := worker.handleResponse(resp, ourApplicationID)
	if err != nil {
		return false, err
	}

	return finished, nil

}

func (worker RabbitMQWorker) sendPollRequest(request models.PollLoanRequest) (*sharedhttp.ClientResponse, error) {
	return worker.httpClient.Get(worker.cfg.BankJobsURL + request.ApplicationID)
}

func (worker RabbitMQWorker) handleResponse(resp *sharedhttp.ClientResponse, ourApplicationID string) (bool, error) {
	switch resp.StatusCode {
	case http.StatusOK:
		return worker.handleOkResponse(resp, ourApplicationID)
	case http.StatusBadRequest:
		return false, errors.New(fmt.Sprintf("Received bad request from bank API : %s", resp.ResponseBody))
	case http.StatusNotFound:
		return false, errors.New(fmt.Sprintf("Received not found response from bank API : %s", resp.ResponseBody))
	default:
		return false, errors.New(fmt.Sprintf("Unexpected response code from bank API : %d : %s", resp.StatusCode, resp.ResponseBody))
	}
}

func (worker RabbitMQWorker) handleOkResponse(resp *sharedhttp.ClientResponse, ourApplicationID string) (bool, error) {
	status, err := getStatusFromResponse(resp)
	if err != nil {
		log.Printf("Could not unmarshal response from bank API %s\n", err)
		return false, err
	}

	if isTerminalStatus(status) {
		// Update the database
		err = worker.repository.UpdateApplicationStatus(ourApplicationID, sharedmodels.Status(status))
		if err != nil {
			log.Printf("Encountered an error updating status in DB %s\n", err)
			return false, err
		}

		fmt.Printf("Marked application %s as %s\n", ourApplicationID, status)
		return true, nil
	}

	// Return false, this tells us to poll again later...
	return false, nil
}

func getStatusFromResponse(response *sharedhttp.ClientResponse) (string, error) {
	var pollLoanSuccessResp *models.PollLoanResponse
	err := json.Unmarshal(response.ResponseBody, &pollLoanSuccessResp)
	if err != nil {
		fmt.Printf("Could not unmarshal poll loan response from bank : %s\n", err)
		return "", err
	}

	return pollLoanSuccessResp.Status, nil
}

func isTerminalStatus(status string) bool {
	return (status == string(sharedmodels.Completed)) || (status == string(sharedmodels.Rejected))
}
