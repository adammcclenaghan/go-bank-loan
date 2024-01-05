package repositorys

import (
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
	"poll-application-service/models"
	"service-shared/http"
	shareddb "service-shared/mocks/database"
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
	wg := &sync.WaitGroup{}
	inChan := make(chan amqp.Delivery)
	cfg := sharedconfig.Config{}
	deliveryHandler := new(sharedmq.DeliveryHandler)
	httpClient := new(sharedhttp.Client)
	repository := new(shareddb.Repository)
	deliveryHandler.On("Nack", false, false, mock.Anything).Return(nil)
	body := "{invalidjson,"

	worker := NewRabbitMQWorker(repository, wg, inChan, cfg, deliveryHandler, httpClient)
	worker.processMessage(getDeliveryWithBody([]byte(body)))

	// Assert that the message is sent to DLQ
	deliveryHandler.AssertCalled(t, "Nack", false, false, mock.Anything)
}

func TestProcessMessageErrorPolling(t *testing.T) {
	delivery := getValidDelivery()
	// Setup
	wg := &sync.WaitGroup{}
	inChan := make(chan amqp.Delivery)
	cfg := sharedconfig.Config{}
	deliveryHandler := new(sharedmq.DeliveryHandler)
	httpClient := new(sharedhttp.Client)
	repository := new(shareddb.Repository)
	deliveryHandler.On("Nack", false, false, delivery).Return(nil)
	httpClient.On("Get", mock.Anything).Return(nil, errors.New(""))

	worker := NewRabbitMQWorker(repository, wg, inChan, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	httpClient.AssertCalled(t, "Get", mock.Anything)
	// Assert that the message is sent to DLQ
	deliveryHandler.AssertCalled(t, "Nack", false, false, delivery)
}

func TestProcessMessageLoanPending(t *testing.T) {
	delivery := getValidDelivery()
	// Setup
	wg := &sync.WaitGroup{}
	inChan := make(chan amqp.Delivery)
	cfg := sharedconfig.Config{}
	deliveryHandler := new(sharedmq.DeliveryHandler)
	httpClient := new(sharedhttp.Client)
	repository := new(shareddb.Repository)
	deliveryHandler.On("Nack", false, true, delivery).Return(nil)
	httpClient.On("Get", mock.Anything).Return(mockLoanStatusResp(string(sharedmodels.Pending)), nil)

	worker := NewRabbitMQWorker(repository, wg, inChan, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	httpClient.AssertCalled(t, "Get", mock.Anything)
	// Assert that the message is re-queued
	deliveryHandler.AssertCalled(t, "Nack", false, true, delivery)
}

func TestProcessMessageLoanCompleted(t *testing.T) {
	delivery := getValidDelivery()
	// Setup
	wg := &sync.WaitGroup{}
	inChan := make(chan amqp.Delivery)
	cfg := sharedconfig.Config{}
	deliveryHandler := new(sharedmq.DeliveryHandler)
	httpClient := new(sharedhttp.Client)
	repository := new(shareddb.Repository)
	deliveryHandler.On("Ack", false, delivery).Return(nil)
	httpClient.On("Get", mock.Anything).Return(mockLoanStatusResp(string(sharedmodels.Completed)), nil)
	repository.On("UpdateApplicationStatus", mock.Anything, mock.Anything).Return(nil)

	worker := NewRabbitMQWorker(repository, wg, inChan, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	httpClient.AssertCalled(t, "Get", mock.Anything)
	// Assert that the message is ack'd
	deliveryHandler.AssertCalled(t, "Ack", false, delivery)
}

func TestProcessMessageLoanRejected(t *testing.T) {
	delivery := getValidDelivery()
	// Setup
	wg := &sync.WaitGroup{}
	inChan := make(chan amqp.Delivery)
	cfg := sharedconfig.Config{}
	deliveryHandler := new(sharedmq.DeliveryHandler)
	httpClient := new(sharedhttp.Client)
	repository := new(shareddb.Repository)
	deliveryHandler.On("Ack", false, delivery).Return(nil)
	httpClient.On("Get", mock.Anything).Return(mockLoanStatusResp(string(sharedmodels.Rejected)), nil)
	repository.On("UpdateApplicationStatus", mock.Anything, mock.Anything).Return(nil)

	worker := NewRabbitMQWorker(repository, wg, inChan, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	httpClient.AssertCalled(t, "Get", mock.Anything)
	// Assert that the message is ack'd
	deliveryHandler.AssertCalled(t, "Ack", false, delivery)
}

func TestProcessMessageInternalDbError(t *testing.T) {
	delivery := getValidDelivery()
	// Setup
	wg := &sync.WaitGroup{}
	inChan := make(chan amqp.Delivery)
	cfg := sharedconfig.Config{}
	deliveryHandler := new(sharedmq.DeliveryHandler)
	httpClient := new(sharedhttp.Client)
	repository := new(shareddb.Repository)
	deliveryHandler.On("Nack", false, false, delivery).Return(nil)
	httpClient.On("Get", mock.Anything).Return(mockLoanStatusResp(string(sharedmodels.Rejected)), nil)
	repository.On("UpdateApplicationStatus", mock.Anything, mock.Anything).Return(errors.New(""))

	worker := NewRabbitMQWorker(repository, wg, inChan, cfg, deliveryHandler, httpClient)
	worker.processMessage(delivery)

	httpClient.AssertCalled(t, "Get", mock.Anything)
	// Assert that the message is sent to DLQ
	deliveryHandler.AssertCalled(t, "Nack", false, false, delivery)
}

func mockLoanStatusResp(status string) *http.ClientResponse {
	return &http.ClientResponse{
		StatusCode:   200,
		ResponseBody: getSuccessResponse(status),
	}
}

func getSuccessResponse(status string) []byte {
	resp := &models.PollLoanResponse{
		ApplicationID: "abc",
		Status:        status,
	}
	bytes, _ := json.Marshal(resp)

	return bytes
}

func getValidDelivery() amqp.Delivery {
	msg := sharedmodels.PollLoanMessage{
		OurApplicationID:  "abc",
		BankApplicationID: "def",
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
