# Dockerfile for chat service

FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o server ./cmd/server

CMD ["./server"]
