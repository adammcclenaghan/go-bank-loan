package shared_config

import "github.com/kelseyhightower/envconfig"

/*
Config defines a set of shared config values shared by services for this assignment.

For example, the following change to the api-gateway can be made in docker-compose to
configure the exposed HTTP port:
api-gateway:
   environment:
      - HTTP_PORT=8082
NB: In the above case a modification to the ports list in docker-compose for the service
is also required.
*/
type Config struct {
	HTTPPort                   int    `envconfig:"http_port" default:"8081"`
	MongoURI                   string `envconfig:"mongo_url" default:"mongodb://application-db:27017"`
	DatabaseName               string `envconfig:"db_name" default:"LoanApplications"`
	DBColletionName            string `envconfig:"db_collection_name" default:"Applications"`
	RabbitMQURL                string `envconfig:"rabbit_mq_url" default:"amqp://guest:guest@rabbit-mq:5672"`
	CreateApplicationQueueName string `envconfig:"create_app_queue_name" default:"create_application"`
	PollApplicationQueueName   string `envconfig:"poll_app_queue_name" default:"poll_applications"`
	PollServiceWorkers         int    `envconfig:"poll_svc_workers" default:"10"`
	CreateServiceWorkers       int    `envconfig:"create_svc_workers" default:"5"`

	BankJobsURL   string `envconfig:"bank_jobs_url" default:"http://bank-api:8000/api/jobs?application_id="`
	BankCreateURL string `envconfig:"bank_create_url" default:"http://bank-api:8000/api/applications"`
}

func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return cfg
}
