package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Booking struct {
	ExpertID        string `json:"expert_id"`
	ScheduledTime   string `json:"scheduled_time"`
	DurationMinutes int    `json:"duration_minutes"`
	Notes           string `json:"notes"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// New booking endpoint
	r.POST("/bookings", func(c *gin.Context) {
		var booking Booking
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// In a real application, you would save the booking to a database here.
		log.Printf("Received new booking: %+v", booking)
		c.JSON(200, gin.H{"message": "Booking received successfully", "booking": booking})
	})

	log.Printf("Booking service is starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
