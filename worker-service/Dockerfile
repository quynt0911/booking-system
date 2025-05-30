FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o worker-service ./cmd/main.go
CMD ["./worker-service"]