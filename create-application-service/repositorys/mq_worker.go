package repositorys

import (
	"bytes"
	"create-application-service/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	sharedhttp "service-shared/http"
	messagequeue "service-shared/message-queue"
	sharedconfig "service-shared/shared-config"
	sharedmodels "service-shared/shared-models"
	"sync"
)

const (
	contentType = "application/json"
)

/*
RabbitMQWorker defines a worker responsible for handling
messages from a rabbitMQ queue. It receives Deliveries via
inChan and will attempt to process them.
*/
type RabbitMQWorker struct {
	wg           *sync.WaitGroup
	inChan       <-chan amqp.Delivery
	publishQueue PublishQueue
	cfg          sharedconfig.Config
	handler      messagequeue.DeliveryHandler
	httpClient   sharedhttp.Client
}

func NewRabbitMQWorker(wg *sync.WaitGroup,
	inChan <-chan amqp.Delivery,
	publishQueue PublishQueue,
	cfg sharedconfig.Config,
	handler messagequeue.DeliveryHandler,
	httpClient sharedhttp.Client) RabbitMQWorker {
	return RabbitMQWorker{
		wg:           wg,
		inChan:       inChan,
		publishQueue: publishQueue,
		cfg:          cfg,
		handler:      handler,
		httpClient:   httpClient,
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
processMesage will process a delivery message.

This will reach out to the create application API of the bank. Once a loan
application has been created with the bank, this worker will delegate responsibility
for publishing a message to the Poll Application Service to a repositorys.PublishQueue

If an unrecoverable error occurs while processing a message, this function will simulate
the behaviour of passing a message to a dead letter queue. Note that the DLQ behaviour has
not yet been implemented due to time constraints.
*/
func (worker RabbitMQWorker) processMessage(delivery amqp.Delivery) {
	var message *sharedmodels.CreateLoanMessage
	err := json.Unmarshal(delivery.Body, &message)
	if messagequeue.CheckError(err,
		fmt.Sprintf("Could not unmarshal message %s to CreateLoanMessage - bad data on queue?", delivery.Body),
		delivery,
		worker.handler) {
		return
	}

	/*
		We generate a new UUID, different from the ApplicationID held in the CreateLoanMessage.
		This is necessary as we may not be the only publisher to the bank API. So to handle cases
		where we encounter a collision in IDs, we should generate them here.
	*/
	bankApplicationID := uuid.New().String()
	loanRequest := models.CreateLoanRequest{
		ID:        bankApplicationID,
		FirstName: message.FirstName,
		LastName:  message.LastName,
	}

	resp, err := worker.sendLoanRequest(loanRequest)
	// Send to DLQ if we cannot contact the bank API. An alternative would be to requeue and try again
	if messagequeue.CheckError(err, "Could not send loan request to bank API", delivery, worker.handler) {
		return
	}

	finished, err := handleCreateResponse(resp)
	// Send to DLQ if we get an unknwon return code from the bank API.
	if messagequeue.CheckError(err, "Unknown return code from bank API", delivery, worker.handler) {
		return
	}

	if !finished {
		// Duplicate UUID - re-queue the msg, we will try again with a new UUID
		worker.handler.Nack(false, true, delivery)
		return
	}

	// We successfully created the loan application
	err = worker.publishQueue.PublishPollRequest(bankApplicationID, message.ApplicationID)
	if messagequeue.CheckError(err, "Created application but could not publish to poll queue", delivery, worker.handler) {
		return
	}

	worker.handler.Ack(false, delivery)
}

func handleCreateResponse(resp *sharedhttp.ClientResponse) (bool, error) {
	switch resp.StatusCode {
	case http.StatusCreated:
		// Successfully created the loan application
		return true, nil
	case http.StatusBadRequest:
		// Duplicate UUID, we must try again
		return false, nil
	default:
		return false, errors.New(fmt.Sprintf("Unexpected response code from bank API %d\n", resp.StatusCode))
	}
}

func (worker RabbitMQWorker) sendLoanRequest(request models.CreateLoanRequest) (*sharedhttp.ClientResponse, error) {
	req, _ := json.Marshal(request)

	fmt.Printf("Sending to URL %s\n", worker.cfg.BankCreateURL)
	return worker.httpClient.Post(worker.cfg.BankCreateURL, contentType, bytes.NewBuffer(req))
}
