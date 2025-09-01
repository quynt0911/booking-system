module github.com/your-org/booking-system/services/worker-service

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/joho/godotenv v1.5.1
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
)