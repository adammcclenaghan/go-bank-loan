//Package repositorys provides interfaces used by controllers to handle data from HTTP requests.
package repositorys

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	sharedconfig "service-shared/shared-config"
	sharedhelpers "service-shared/shared-helpers"
	sharedmodels "service-shared/shared-models"
)

//PublishQueue defines an interface for interacting with a message queue. Specifically, it provides
//an abstraction for sending a loan request to a message queue.
type PublishQueue interface {
	PublishLoanRequest(loanApplication sharedmodels.CreateLoanMessage) error
}

//RabbitMessageQueue is used to interact with a RabbitMQ message queue.
type RabbitMessageQueue struct {
	queue *amqp.Queue
	ch    *amqp.Channel
	cfg   sharedconfig.Config
}

//NewRabbitQueue returns a RabbitMessageQueue struct.
func NewRabbitQueue(ch *amqp.Channel, cfg sharedconfig.Config) *RabbitMessageQueue {
	queue, err := ch.QueueDeclare(
		cfg.CreateApplicationQueueName, // name
		true,                           // durable (survive restarts)
		false,                          // do not delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil)                            // args

	sharedhelpers.FailOnError(err, "Gateway publisher failed to declare create application queue")
	return &RabbitMessageQueue{queue: &queue, ch: ch}
}

/*
PublishLoanRequest publishes a message to a RabbitMQ queue. This message is intended to be consumed
by a consumer, which should then negotiate with the bank API and create a loan application.
*/
func (msgQueue RabbitMessageQueue) PublishLoanRequest(createRequest sharedmodels.CreateLoanMessage) error {
	fmt.Printf("Publishing loan request %#v\n", createRequest)
	request, _ := json.Marshal(createRequest)

	publishErr := msgQueue.ch.Publish(
		"",
		msgQueue.queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         request,
		})

	return publishErr
}
