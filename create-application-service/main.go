package main

import (
	"create-application-service/repositorys"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	sharedhttp "service-shared/http"
	messagequeue "service-shared/message-queue"
	sharedconfig "service-shared/shared-config"
	sharedhelpers "service-shared/shared-helpers"
	"sync"
)

/*
main Starts the Create Application Service

This service is a consumer service. It is responsible for retrieving
messages off a queue. These messages instruct this service to create a
new loan application with the bank API.

This consumer also acts as a publisher, after creating an application with
the bank API, it will publish a message to a queue. This message is intended
to be picked up by the Poll Application Service.

This function spawns parallel workers.
For implementation details, please see repositorys.mq_worker
*/
func main() {
	cfg := sharedconfig.Get()

	// Connect to the rabbit queue that we will consume from
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	sharedhelpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create publish queue
	ch, err := conn.Channel()
	sharedhelpers.FailOnError(err, "Publisher failed to open a channel to RabbitMQ")
	defer ch.Close()
	publishQueue := repositorys.NewRabbitPublishQueue(ch, cfg)

	// Set up worker to consume off the channel and publish to the poll queue
	in := make(chan amqp.Delivery)
	wg := &sync.WaitGroup{}
	maxWorkers := cfg.CreateServiceWorkers
	wg.Add(cfg.CreateServiceWorkers)
	handler := messagequeue.AmqpDeliveryHandler{}
	httpClient := sharedhttp.DefaultClient{HttpClient: http.DefaultClient}
	for i := 1; i <= maxWorkers; i++ {
		worker := repositorys.NewRabbitMQWorker(wg, in, publishQueue, cfg, handler, httpClient)
		go worker.ProcessMessages()
	}

	// Consumes messages from the queue, passes to in, which is consumed by the workers
	consumer := messagequeue.NewRabbitMQConsumer(conn, maxWorkers, in, cfg.CreateApplicationQueueName)
	go consumer.Consume()

	wg.Wait()
}
