package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Booking struct {
	ID              uuid.UUID  `db:"id"`
	UserID          uuid.UUID  `db:"user_id"`
	ExpertID        uuid.UUID  `db:"expert_id"`
	ScheduledTime   time.Time  `db:"scheduled_datetime"`
	DurationMinutes int        `db:"duration_minutes"`
	MeetingType     string     `db:"meeting_type"`
	MeetingURL      *string    `db:"meeting_url"`
	MeetingAddress  *string    `db:"meeting_address"`
	Notes           string     `db:"notes"`
	Status          string     `db:"status"`
	Price           *float64   `db:"price"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	ConfirmedAt     *time.Time `db:"confirmed_at"`
	CancelledAt     *time.Time `db:"cancelled_at"`
	CompletedAt     *time.Time `db:"completed_at"`
}

type BookingRequest struct {
	UserID    string `json:"user_id"`
	ExpertID  string `json:"expert_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	dsn := os.Getenv("BOOKING_SERVICE_DSN")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/CreateBooking", func(c *gin.Context) {
		var req BookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		scheduledTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}

		endTime, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}

		duration := int(endTime.Sub(scheduledTime).Minutes())

		booking := Booking{
			ID:              uuid.New(),
			UserID:          uuid.MustParse(req.UserID),
			ExpertID:        uuid.MustParse(req.ExpertID),
			ScheduledTime:   scheduledTime,
			DurationMinutes: duration,
			MeetingType:     "online",
			MeetingURL:      nil,
			MeetingAddress:  nil,
			Notes:           req.Notes,
			Status:          req.Status,
			Price:           nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ConfirmedAt:     nil,
			CancelledAt:     nil,
			CompletedAt:     nil,
		}

		_, err = db.NamedExec(`INSERT INTO bookings 
			(id, user_id, expert_id, scheduled_datetime, duration_minutes, meeting_type, meeting_url, meeting_address, notes, status, price, created_at, updated_at, confirmed_at, cancelled_at, completed_at)
			VALUES (:id, :user_id, :expert_id, :scheduled_datetime, :duration_minutes, :meeting_type, :meeting_url, :meeting_address, :notes, :status, :price, :created_at, :updated_at, :confirmed_at, :cancelled_at, :completed_at)`, &booking)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save booking: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Booking created successfully", "booking": booking})
	})

	log.Printf("Booking service is starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
