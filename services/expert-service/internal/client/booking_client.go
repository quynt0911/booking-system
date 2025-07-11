package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Booking struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	ExpertID        string    `json:"expert_id"`
	ScheduledTime   time.Time `json:"scheduled_datetime"`
	DurationMinutes int       `json:"duration_minutes"`
	EndTime         time.Time
	Status          string `json:"status"`
}

type BookingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Bookings []Booking `json:"bookings"`
	} `json:"data"`
}

// GetBookingsByExpertAndDateRange gọi API booking-service để lấy danh sách booking của expert trong khoảng ngày
func GetBookingsByExpertAndDateRange(bookingServiceURL, expertID, token string, startDate, endDate time.Time) ([]Booking, error) {
	url := fmt.Sprintf("%s/GetExpertBookings?expert_id=%s&start_date=%s&end_date=%s",
		bookingServiceURL, expertID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Log để debug
	log.Printf("Calling booking service: %s", url)
	log.Printf("Token: %s", token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set Authorization header
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	} else {
		log.Printf("Warning: No token provided for booking service request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling booking service: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Log response status
	log.Printf("Booking service response status: %d", resp.StatusCode)

	// Read and log raw response for debugging
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Booking service raw response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("booking service returned status %d: %s", resp.StatusCode, string(body))
	}

	var result BookingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error decoding booking service response: %v", err)
		return nil, err
	}

	// Calculate end time for each booking
	var bookings []Booking
	for _, b := range result.Data.Bookings {
		if b.Status == "pending" || b.Status == "confirmed" {
			b.EndTime = b.ScheduledTime.Add(time.Duration(b.DurationMinutes) * time.Minute)
			bookings = append(bookings, b)
		}
	}

	log.Printf("Successfully retrieved %d bookings", len(bookings))
	for _, b := range bookings {
		log.Printf("Booking: time=%s-%s, status=%s",
			b.ScheduledTime.Format("2006-01-02 15:04:05"),
			b.EndTime.Format("2006-01-02 15:04:05"),
			b.Status)
	}

	return bookings, nil
}
