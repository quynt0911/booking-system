package main

import (
	"database/sql"
	"expert-service/internal/cache"
	"expert-service/internal/config"
	"expert-service/internal/handler"
	"expert-service/internal/repository"
	"expert-service/internal/routes"
	"expert-service/internal/service"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to PostgreSQL
	db, err := connectDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db, cfg.Database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Connect to Redis
	redisClient := connectRedis(cfg.Redis)
	defer redisClient.Close()

	// Initialize repositories
	expertRepo := repository.NewExpertRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	offTimeRepo := repository.NewOffTimeRepository(db)

	// Initialize cache
	availabilityCache := cache.NewAvailabilityCache(redisClient, cfg.Cache.TTL)

	// Initialize services
	expertSvc := service.NewExpertService(expertRepo)
	scheduleSvc := service.NewScheduleService(scheduleRepo)
	availabilitySvc := service.NewAvailabilityService(scheduleRepo, offTimeRepo, availabilityCache)

	// Initialize handlers
	expertHandler := handler.NewExpertHandler(expertSvc)
	scheduleHandler := handler.NewScheduleHandler(scheduleSvc)
	availabilityHandler := handler.NewAvailabilityHandler(availabilitySvc)

	// Setup routes
	router := gin.Default()
	
	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Setup API routes
	routes.SetupRoutes(router, expertHandler, scheduleHandler, availabilityHandler)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"service":   "expert-service",
		})
	})

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting expert-service on %s", serverAddr)
	
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}

func connectRedis(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis")
	return client
}

func runMigrations(db *sql.DB, cfg config.DatabaseConfig) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}