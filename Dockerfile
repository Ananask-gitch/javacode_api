FROM golang:1.23.6-alpine3.21
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd

