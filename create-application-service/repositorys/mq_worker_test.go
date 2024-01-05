package repositorys

import (
	mocks "create-application-service/mocks/repositorys"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
	sharedhttp2 "service-shared/http"
	sharedhttp "service-shared/mocks/http"
	sharedmq "service-shared/mocks/message-queue"
	sharedconfig "service-shared/shared-config"
	sharedmodels "service-shared/shared-models"
	"sync"
	"testing"
	"time"
)

func TestProcessMessageFailsToUnmarshalBody(t *testing.T) {
	// Setup
	publishQueue := new(mocks.PublishQueue)
	wg := &sync.WaitGroup{}
	cfg := sharedconfig.Config{}
	inChan := make(chan amqp.Delivery)
	deliveryHandler := new(sharedmq.DeliveryHandler)
	deliveryHandler.On("Nack", false, false, mock.Anything).Return(nil)

	// Create worker
	worker := NewRabbitMQWorker(wg, inChan, publishQueue, cfg, deliveryHandler, new(sharedhttp.Client))

	body := "{invalidjson,"
	worker.processMessage(getDeliveryWithBody([]byte(body)))

	// Assert that the message is sent to DLQ
	deliveryHandler.AssertCalled(t, "Nack", false, false, mock.Anything)
}

func TestProcessMessageFailsToSendRequestToBank(t *testing.T) {
	delivery := getValidDelivery()
	// Setup
	publishQueue := new(mocks.PublishQueue)
	wg := &sync.WaitGroup{}
	cfg := sharedconfig.Config{}
	inChan := make(chan amqp.Delivery)
	deliveryHandler := new(sharedmq.DeliveryHandler)
	deliveryHandler.On("Nack", false, false, delivery).Return(nil)
	httpClient := new(sharedhttp.Client)
	httpClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New(""))

	// Create worker
	worker := NewRabbitMQWorker(wg, inChan, publishQueue, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	// Assert that the message is sent to the dlq
	deliveryHandler.AssertCalled(t, "Nack", false, false, delivery)
}

func TestProcessMessageUnknownBankStatusCode(t *testing.T) {
	delivery := getValidDelivery()
	// Setup
	publishQueue := new(mocks.PublishQueue)
	wg := &sync.WaitGroup{}
	cfg := sharedconfig.Config{}
	inChan := make(chan amqp.Delivery)
	deliveryHandler := new(sharedmq.DeliveryHandler)
	deliveryHandler.On("Nack", false, false, delivery).Return(nil)
	httpClient := new(sharedhttp.Client)
	response := &sharedhttp2.ClientResponse{
		StatusCode:   1,
		ResponseBody: nil,
	}
	httpClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(response, nil)

	// Create worker
	worker := NewRabbitMQWorker(wg, inChan, publishQueue, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	// Assert that the message is sent to the dlq
	deliveryHandler.AssertCalled(t, "Nack", false, false, delivery)
}

func TestProcessMessageDuplicateIdStatusCode(t *testing.T) {
	delivery := getValidDelivery()
	// Setup
	publishQueue := new(mocks.PublishQueue)
	wg := &sync.WaitGroup{}
	cfg := sharedconfig.Config{}
	inChan := make(chan amqp.Delivery)
	deliveryHandler := new(sharedmq.DeliveryHandler)
	deliveryHandler.On("Nack", false, true, delivery).Return(nil)
	httpClient := new(sharedhttp.Client)
	response := &sharedhttp2.ClientResponse{
		StatusCode:   400,
		ResponseBody: nil,
	}
	httpClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(response, nil)

	// Create worker
	worker := NewRabbitMQWorker(wg, inChan, publishQueue, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	// Assert that the message is requeued
	deliveryHandler.AssertCalled(t, "Nack", false, true, delivery)
}

func TestProcessMessageFailToPublishToPollQueue(t *testing.T) {
	delivery := getValidDelivery()

	// Setup
	publishQueue := new(mocks.PublishQueue)
	wg := &sync.WaitGroup{}
	cfg := sharedconfig.Config{}
	inChan := make(chan amqp.Delivery)
	deliveryHandler := new(sharedmq.DeliveryHandler)
	deliveryHandler.On("Nack", false, false, delivery).Return(nil)
	publishQueue.On("PublishPollRequest", mock.Anything, mock.Anything).Return(errors.New(""))
	httpClient := new(sharedhttp.Client)
	response := &sharedhttp2.ClientResponse{
		StatusCode:   201,
		ResponseBody: nil,
	}
	httpClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(response, nil)

	// Create worker
	worker := NewRabbitMQWorker(wg, inChan, publishQueue, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	// Assert that the message is sent to DLQ
	publishQueue.AssertCalled(t, "PublishPollRequest", mock.Anything, mock.Anything)
	deliveryHandler.AssertCalled(t, "Nack", false, false, delivery)
}

func TestProcessMessageAckWhenSuccess(t *testing.T) {
	delivery := getValidDelivery()

	// Setup
	publishQueue := new(mocks.PublishQueue)
	wg := &sync.WaitGroup{}
	cfg := sharedconfig.Config{}
	inChan := make(chan amqp.Delivery)
	deliveryHandler := new(sharedmq.DeliveryHandler)
	deliveryHandler.On("Ack", false, delivery).Return(nil)
	publishQueue.On("PublishPollRequest", mock.Anything, mock.Anything).Return(nil)
	httpClient := new(sharedhttp.Client)
	response := &sharedhttp2.ClientResponse{
		StatusCode:   201,
		ResponseBody: nil,
	}
	httpClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(response, nil)

	// Create worker
	worker := NewRabbitMQWorker(wg, inChan, publishQueue, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	// Assert that the message is sent to DLQ
	publishQueue.AssertCalled(t, "PublishPollRequest", mock.Anything, mock.Anything)
	// Assert the message is ACKd
	deliveryHandler.AssertCalled(t, "Ack", false, delivery)
}

func getValidDelivery() amqp.Delivery {
	msg := sharedmodels.CreateLoanMessage{
		ApplicationID: "Test",
		FirstName:     "First",
		LastName:      "Last",
	}
	bytes, _ := json.Marshal(msg)

	return getDeliveryWithBody(bytes)
}

func getDeliveryWithBody(body []byte) amqp.Delivery {
	return amqp.Delivery{
		Acknowledger:    nil,
		Headers:         nil,
		ContentType:     "",
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		ConsumerTag:     "",
		MessageCount:    0,
		DeliveryTag:     0,
		Redelivered:     false,
		Exchange:        "",
		RoutingKey:      "",
		Body:            body,
	}
}
