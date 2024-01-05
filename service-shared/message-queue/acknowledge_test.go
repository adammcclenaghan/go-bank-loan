package message_queue

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	mocks "service-shared/mocks/message-queue"
	"testing"
)

func TestCheckErrorWhenError(t *testing.T) {
	err := errors.New("")
	msg := ""
	delivery := amqp.Delivery{}
	handler := new(mocks.DeliveryHandler)
	handler.On("Nack", false, false, delivery).Return(nil)

	hadError := CheckError(err, msg, delivery, handler)

	assert.True(t, hadError)
	handler.AssertCalled(t, "Nack", false, false, delivery)
}

func TestCheckErrorWhenNoError(t *testing.T) {
	var err error = nil
	msg := ""
	delivery := amqp.Delivery{}
	handler := new(mocks.DeliveryHandler)
	handler.On("Nack", false, delivery).Return(nil)

	hadError := CheckError(err, msg, delivery, handler)

	assert.False(t, hadError)
	handler.AssertNotCalled(t, "Nack", false, delivery)
}
