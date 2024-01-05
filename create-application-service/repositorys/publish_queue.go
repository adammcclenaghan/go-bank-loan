package repositorys

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	sharedconfig "service-shared/shared-config"
	sharedhelpers "service-shared/shared-helpers"
	sharedmodels "service-shared/shared-models"
)

//PublishQueue defines an interface for interacting with a message queue. Specifically, it provies
//an abstraction for sending a poll loan request to a message queue.
type PublishQueue interface {
	PublishPollRequest(bankApplicationID, ourApplicationID string) error
}

//RabbitPublishQueue is used to interact with a RabbitMQ message queue.
type RabbitPublishQueue struct {
	queue *amqp.Queue
	ch    *amqp.Channel
	cfg   sharedconfig.Config
}

//NewRabbitPublishQueue returns a RabbitPublishQueue struct.
func NewRabbitPublishQueue(ch *amqp.Channel, cfg sharedconfig.Config) *RabbitPublishQueue {
	pollQueue, err := ch.QueueDeclare(
		cfg.PollApplicationQueueName, // name
		true,                         // durable (survive restarts)
		false,                        // do not delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil)                          // args

	sharedhelpers.FailOnError(err, "Publisher failed to declare the poll application queue")
	return &RabbitPublishQueue{queue: &pollQueue, ch: ch}
}

/*
PublishPollRequest publishes a message to a RabbitMQ queue. This message is intended to be consumed by
a consumer, which should then negotiate with the jobs API of the bank to determine the status of an application.
*/
func (queue RabbitPublishQueue) PublishPollRequest(bankApplicationID, ourApplicationID string) error {
	message := sharedmodels.PollLoanMessage{BankApplicationID: bankApplicationID, OurApplicationID: ourApplicationID}
	request, _ := json.Marshal(message)

	publishErr := queue.ch.Publish(
		"",
		queue.queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         request,
		})

	return publishErr
}
