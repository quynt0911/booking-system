if [ -z "$1" ]; then
    echo "Usage: ./dev-shell.sh <service-name>"
    echo "Available services: user-service, booking-service, expert-service, api-gateway"
    exit 1
fi

SERVICE_NAME=$1
docker-compose -f docker-compose.dev.yml exec $SERVICE_NAME sh