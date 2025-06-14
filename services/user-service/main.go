package main

import (
	"log"
	"os"
	"services/user-service/handler"
	"services/user-service/repository"
	"services/user-service/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("USER_SERVICE_DSN")
	if dsn == "" {
		log.Fatal("USER_SERVICE_DSN is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// Remove auto-migrate since we're using our own schema
	// db.AutoMigrate(&model.User{})

	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiry := 3600
	refreshExpiry := 604800
	jwtService := service.NewJWTService(jwtSecret, jwtExpiry, refreshExpiry)

	r := gin.Default()
	handler.RegisterRoutes(r, userService, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
