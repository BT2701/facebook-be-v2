# Sử dụng base image Golang
FROM golang:1.22

WORKDIR /app

# Copy code và module Go
COPY . .

# Tải các dependency
RUN go mod tidy

# Build binary
RUN go build -o notification-service

# Chạy ứng dụng
CMD ["./notification-service"]
