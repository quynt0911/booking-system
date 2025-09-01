FROM golang:1.24.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o expert-service ./cmd/main.go

EXPOSE 8084

CMD ["./notification-service"]
