module create-application-service

go 1.18

require service-shared v0.0.0

replace service-shared v0.0.0 => ../service-shared

require (
	github.com/google/uuid v1.3.0
	github.com/rabbitmq/amqp091-go v1.3.4
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	go.mongodb.org/mongo-driver v1.9.1 // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3 // indirect
	golang.org/x/sys v0.0.0-20210806184541-e5e7981a1069 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
