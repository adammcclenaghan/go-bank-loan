# Alpine image as it is small
FROM golang:1.18-alpine

WORKDIR /app/service-shared
ADD service-shared .

WORKDIR /app/api-gateway
ADD api-gateway .

RUN go mod download

RUN go build -o /app/api-gateway/main

CMD ["/app/api-gateway/main"]


