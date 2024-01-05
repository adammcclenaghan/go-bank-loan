package message_queue

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

//DeliveryHandler provides an abstraction for handling RabbitMQ deliveries.
//It allows for easier testing via an interface which unit tests can mock.
type DeliveryHandler interface {
	Ack(multiple bool, delivery amqp.Delivery) error
	Nack(multiple bool, requeue bool, delivery amqp.Delivery) error
	Reject(requeue bool, delivery amqp.Delivery) error
}

type AmqpDeliveryHandler struct{}

func (handler AmqpDeliveryHandler) Ack(multiple bool, delivery amqp.Delivery) error {
	return delivery.Ack(multiple)
}

func (handler AmqpDeliveryHandler) Nack(multiple bool, requeue bool, delivery amqp.Delivery) error {
	return delivery.Nack(multiple, requeue)
}

func (handler AmqpDeliveryHandler) Reject(requeue bool, delivery amqp.Delivery) error {
	return delivery.Reject(requeue)
}

/*
CheckError is a helper function. In the event that an error occurs, it will
print the provided message and send the amqp delivery to the DLQ
Returns true iff err is not nil.
*/
func CheckError(err error, msg string, delivery amqp.Delivery, handler DeliveryHandler) bool {
	if err != nil {
		fmt.Printf("%s : %s\n", msg, err)
		// Future improvement : Setup queue so that Nack with requeue=false goes to DLQ
		handler.Nack(false, false, delivery)
		return true
	}

	return false
}
