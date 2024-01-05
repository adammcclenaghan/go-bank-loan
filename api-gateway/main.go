package main

import (
	"api-gateway/controllers"
	"api-gateway/docs"
	"api-gateway/repositorys"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"io"
	"log"
	"os"
	"service-shared/database"
	shared_config "service-shared/shared-config"
	sharedhelpers "service-shared/shared-helpers"
	sharedmodels "service-shared/shared-models"
)

/*
main Starts the API gateway.

This API gateway is capable of service the following requests:
1. Creating a loan application
2. Getting the status of a loan application, given its application ID
3. Getting all applications with a given status.

This application is responsible for receiving and handling HTTP requests.

GET requests do not reach out to the bank API, but provide a view of loan applications
directly from the database to clients.

POST requests also do not reach out to the bank API, instead they publish a message to a message
queue, and return the status of an application as pending should the create request be successful.
This means that this API gateway does not rely on an 'live' bank API to provide some response to clients.
*/
func main() {
	cfg := shared_config.Get()
	fmt.Println("API Gateway is starting ...")

	// Connect to database
	fmt.Println("Connecting to db ... ")
	dbClient, err := database.InitClient(cfg)
	sharedhelpers.FailOnError(err, "Failed to initialise the db client")
	defer dbClient.Disconnect(context.Background())
	collection := dbClient.Database(cfg.DatabaseName).Collection(cfg.DBColletionName)
	database.InitIndexes(collection)
	mongo := database.NewMongoCollection(collection)
	repository := database.NewMongoRepository(mongo)

	// Setup rabbitmq work queue
	fmt.Println("Connecting to RabbitMQ ... ")
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	sharedhelpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	sharedhelpers.FailOnError(err, "Failed to open a channel to RabbitMQ")
	messageQueue := repositorys.NewRabbitQueue(ch, cfg)

	// Set up controller
	controller := controllers.NewLoanAppController(repository, messageQueue)

	// Setup the API
	fmt.Println("Setting up the API router ...")
	// Disable console color as we are going to write gin logs to file
	gin.DisableConsoleColor()
	// Tell gin to log to a file
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Add custom validator for validstatus
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("validstatus", sharedmodels.ValidStatus)
		sharedhelpers.FailOnError(err, "Failed to registry the validstatus validator")
	}

	router := gin.Default()
	router.POST("/api/application", controller.CreateApplication)
	router.GET("/api/application", controller.GetApplication)
	router.GET("/api/applications-with-status", controller.GetApplicationsWithStatus)
	docs.SwaggerInfo.Title = "Go Bank Loan API"
	docs.SwaggerInfo.Description = "An API which simulates creating loans with a banking API, as well as receiving information about the status of those loans."
	docs.SwaggerInfo.Version = "1.0"
	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("API is UP ... ")

	err = router.Run(fmt.Sprintf(":%d", cfg.HTTPPort))
	if err != nil {
		log.Fatal(err)
	}
}
