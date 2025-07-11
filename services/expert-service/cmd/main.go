package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"expert-service/internal/cache"
	"expert-service/internal/repository"
	"expert-service/internal/routes"
	"expert-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Lấy DSN từ biến môi trường
	dsn := os.Getenv("EXPERT_SERVICE_DSN")
	if dsn == "" {
		log.Fatal("EXPERT_SERVICE_DSN is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
	})
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	availabilityCache := cache.NewAvailabilityCache(redisClient, time.Hour)

	// Initialize repositories
	expertRepo := repository.NewExpertRepository(db)
	expertScheduleRepo := repository.NewExpertScheduleRepository(db)
	offTimeRepo := repository.NewOffTimeRepository(db)

	// Initialize services
	expertService := service.NewExpertService(expertRepo)
	expertScheduleService := service.NewExpertScheduleService(expertScheduleRepo)
	availabilityService := service.NewExpertAvailabilityService(expertRepo, repository.NewScheduleRepository(db), offTimeRepo, expertScheduleRepo, availabilityCache)

	// Initialize router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, expertService, expertScheduleService, availabilityService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	router.Run(":" + port)
}
