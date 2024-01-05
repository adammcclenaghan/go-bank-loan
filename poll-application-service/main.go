package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"poll-application-service/repositorys"
	"service-shared/database"
	sharedhttp "service-shared/http"
	messagequeue "service-shared/message-queue"
	sharedconfig "service-shared/shared-config"
	helpers "service-shared/shared-helpers"
	"sync"
)

/*
main Starts the Poll Application Service

This service is a consumer service. It is responsible for retrieving
messages off a queue. These messages instruct this service to poll the
bank API for a given loan application.

This function spawns parallel workers.
For implementation details, please see repositorys.mq_worker.
*/
func main() {
	cfg := sharedconfig.Get()

	fmt.Println("Connecting to db ... ")
	dbClient, err := database.InitClient(cfg)
	helpers.FailOnError(err, "Failed to initialise the db client")
	defer dbClient.Disconnect(context.Background())
	collection := dbClient.Database(cfg.DatabaseName).Collection(cfg.DBColletionName)
	database.InitIndexes(collection)
	mongo := database.NewMongoCollection(collection)
	repository := database.NewMongoRepository(mongo)

	fmt.Println("Connecting to RabbitMQ ... ")
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Workers to process messages received from the queue
	fmt.Println("Creating workers to consume from rabbit MQ")
	in := make(chan amqp.Delivery)
	wg := &sync.WaitGroup{}
	maxWorkers := cfg.PollServiceWorkers
	handler := messagequeue.AmqpDeliveryHandler{}
	httpClient := sharedhttp.DefaultClient{HttpClient: http.DefaultClient}
	wg.Add(maxWorkers)
	for i := 1; i <= maxWorkers; i++ {
		worker := repositorys.NewRabbitMQWorker(repository, wg, in, cfg, handler, httpClient)
		go worker.ProcessMessages()
	}
	// Consumes messages from the queue, passes to in, which is consumed by the workers
	consumer := messagequeue.NewRabbitMQConsumer(conn, maxWorkers, in, cfg.PollApplicationQueueName)

	go consumer.Consume()
	wg.Wait()
}
