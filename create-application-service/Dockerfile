# Alpine image as it is small
FROM golang:1.18-alpine

WORKDIR /app/service-shared
ADD service-shared .

WORKDIR /app/create-application-service
ADD create-application-service .


RUN go mod download

RUN go build -o /app/create-application-service/main

CMD ["/app/create-application-service/main"]


