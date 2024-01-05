# Alpine image as it is small
FROM golang:1.21

WORKDIR /app/bank-api
ADD bank-api .

WORKDIR /app/bank-api
ADD bank-api .

RUN go mod download

RUN go build -o /app/bank-api/main

CMD ["/app/bank-api/main"]

