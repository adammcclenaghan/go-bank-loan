package message_queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	sharedhelpers "service-shared/shared-helpers"
)

//Consumer provides an interface for consumer services
type Consumer interface {
	Consume()
}

/*
RabbitMQConsumer represents a type capable of consuming off a RabbitMQ queue.
Currently, only the QueueName and PrefetchSize are configurable via this struct.
Of course, this could be extended in future to allow configuring of the queue or
channel args etc.
*/
type RabbitMQConsumer struct {
	connection   *amqp.Connection
	prefetchSize int
	outChan      chan<- amqp.Delivery
	queueName    string
}

//NewRabbitMQConsumer creates a RabbitMQConsumer
func NewRabbitMQConsumer(conn *amqp.Connection, prefetchSize int, outChan chan<- amqp.Delivery, queueName string) RabbitMQConsumer {
	return RabbitMQConsumer{
		connection:   conn,
		prefetchSize: prefetchSize,
		outChan:      outChan,
		queueName:    queueName,
	}
}

/*
Consume sets up and iterates over messages from the channel.
Each message is delivered to RabbitMQConsumer.OutChan

The intent is that there should be workers listening on this channel. Those
workers are responsible for processing the contents of a delivered message.
*/
func (consumer RabbitMQConsumer) Consume() {
	ch, err := consumer.connection.Channel()
	sharedhelpers.FailOnError(err, "Consumer failed to open a channel to RabbitMQ")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		consumer.queueName, // name
		true,               // durable (survive restarts)
		false,              // do not delete when unused
		false,              // exclusive
		false,              // no-wait
		nil)                /// args
	sharedhelpers.FailOnError(err, "Failed to declare the create application queue")

	// https://www.rabbitmq.com/consumer-prefetch.html
	err = ch.Qos(
		consumer.prefetchSize, // prefetch count
		0,                     // prefetch size
		false)                 // global
	sharedhelpers.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consume
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      //no-wait
		nil)        // args
	sharedhelpers.FailOnError(err, "Failed to register consumer")

	var forever chan struct{}
	go func() {
		for d := range msgs {
			consumer.outChan <- d
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}
