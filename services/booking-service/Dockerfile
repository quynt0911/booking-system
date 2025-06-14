FROM golang:1.24.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o booking-service ./cmd/main.go

EXPOSE 8082

CMD ["./booking-service"]