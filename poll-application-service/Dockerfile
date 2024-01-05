# Alpine image as it is small
FROM golang:1.18-alpine

WORKDIR /app/service-shared
ADD service-shared .

WORKDIR /app/poll-application-service
ADD poll-application-service .

RUN go mod download

RUN go build -o /app/poll-application-service/main

CMD ["/app/poll-application-service/main"]


