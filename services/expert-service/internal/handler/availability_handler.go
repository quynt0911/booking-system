package handler

import (
	"encoding/json"
	"expert-service/internal/model"
	"expert-service/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

// AvailabilityHandler handles availability-related requests
type AvailabilityHandler struct {
	availabilityService service.ExpertAvailabilityService
}

// NewAvailabilityHandler creates a new availability handler
func NewAvailabilityHandler(availabilityService service.ExpertAvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{
		availabilityService: availabilityService,
	}
}

// Availability represents expert availability
type Availability struct {
	ID        int       `json:"id"`
	ExpertID  int       `json:"expert_id"`
	Date      string    `json:"date"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	IsBooked  bool      `json:"is_booked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateAvailabilityRequest represents request body for creating availability
type CreateAvailabilityRequest struct {
	ExpertID  int    `json:"expert_id" validate:"required"`
	Date      string `json:"date" validate:"required"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

// UpdateAvailabilityRequest represents request body for updating availability
type UpdateAvailabilityRequest struct {
	Date      string `json:"date,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	IsBooked  *bool  `json:"is_booked,omitempty"`
}

// CreateAvailability creates new availability slot
func (h *AvailabilityHandler) CreateAvailability(w http.ResponseWriter, r *http.Request) {
	var req CreateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate time format
	if !isValidTimeFormat(req.StartTime) || !isValidTimeFormat(req.EndTime) {
		http.Error(w, "Invalid time format. Use HH:MM", http.StatusBadRequest)
		return
	}

	// Validate date format
	if !isValidDateFormat(req.Date) {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to create availability
	// Example response
	availability := Availability{
		ID:        1,
		ExpertID:  req.ExpertID,
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IsBooked:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    availability,
		"message": "Availability created successfully",
	})
}

// GetAvailabilities retrieves all availabilities with optional filters
func (h *AvailabilityHandler) GetAvailabilities(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	expertID := r.URL.Query().Get("expert_id")
	date := r.URL.Query().Get("date")
	isBooked := r.URL.Query().Get("is_booked")

	// TODO: Add business logic to fetch availabilities with filters
	// Example response
	availabilities := []Availability{
		{
			ID:        1,
			ExpertID:  1,
			Date:      "2024-12-10",
			StartTime: "09:00",
			EndTime:   "10:00",
			IsBooked:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			ExpertID:  1,
			Date:      "2024-12-10",
			StartTime: "14:00",
			EndTime:   "15:00",
			IsBooked:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    availabilities,
		"filters": map[string]string{
			"expert_id": expertID,
			"date":      date,
			"is_booked": isBooked,
		},
	})
}

// GetAvailabilityByID retrieves availability by ID
func (h *AvailabilityHandler) GetAvailabilityByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid availability ID", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to fetch availability by ID
	// Example response
	availability := Availability{
		ID:        id,
		ExpertID:  1,
		Date:      "2024-12-10",
		StartTime: "09:00",
		EndTime:   "10:00",
		IsBooked:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    availability,
	})
}

// UpdateAvailability updates existing availability
func (h *AvailabilityHandler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid availability ID", http.StatusBadRequest)
		return
	}

	var req UpdateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate time format if provided
	if req.StartTime != "" && !isValidTimeFormat(req.StartTime) {
		http.Error(w, "Invalid start time format. Use HH:MM", http.StatusBadRequest)
		return
	}

	if req.EndTime != "" && !isValidTimeFormat(req.EndTime) {
		http.Error(w, "Invalid end time format. Use HH:MM", http.StatusBadRequest)
		return
	}

	// Validate date format if provided
	if req.Date != "" && !isValidDateFormat(req.Date) {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to update availability
	// Example response
	availability := Availability{
		ID:        id,
		ExpertID:  1,
		Date:      getValueOrDefault(req.Date, "2024-12-10"),
		StartTime: getValueOrDefault(req.StartTime, "09:00"),
		EndTime:   getValueOrDefault(req.EndTime, "10:00"),
		IsBooked:  getBoolValueOrDefault(req.IsBooked, false),
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    availability,
		"message": "Availability updated successfully",
	})
}

// DeleteAvailability deletes availability by ID
func (h *AvailabilityHandler) DeleteAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid availability ID", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to delete availability
	// Check if availability exists and is not booked before deletion

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Availability deleted successfully",
		"id":      id,
	})
}

// GetExpertAvailabilities retrieves all availabilities for a specific expert
func (h *AvailabilityHandler) GetExpertAvailabilities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	expertID, err := strconv.Atoi(vars["expert_id"])
	if err != nil {
		http.Error(w, "Invalid expert ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters for additional filters
	date := r.URL.Query().Get("date")
	isBooked := r.URL.Query().Get("is_booked")

	// TODO: Add business logic to fetch expert's availabilities
	// Example response
	availabilities := []Availability{
		{
			ID:        1,
			ExpertID:  expertID,
			Date:      "2024-12-10",
			StartTime: "09:00",
			EndTime:   "10:00",
			IsBooked:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			ExpertID:  expertID,
			Date:      "2024-12-10",
			StartTime: "14:00",
			EndTime:   "15:00",
			IsBooked:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"data":      availabilities,
		"expert_id": expertID,
		"filters": map[string]string{
			"date":      date,
			"is_booked": isBooked,
		},
	})
}

// BookAvailability books an availability slot
func (h *AvailabilityHandler) BookAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid availability ID", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to book availability
	// Check if availability is not already booked
	// Update the is_booked status

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Availability booked successfully",
		"id":      id,
	})
}

// CheckAvailability handles availability checks
func (h *AvailabilityHandler) CheckAvailability(c *gin.Context) {
	var req model.CheckAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	isAvailable, err := h.availabilityService.CheckAvailability(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"is_available": isAvailable,
		},
	})
}

// Helper functions
func isValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

func isValidDateFormat(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func getValueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getBoolValueOrDefault(value *bool, defaultValue bool) bool {
	if value == nil {
		return defaultValue
	}
	return *value
}

func (h *AvailabilityHandler) CreateOffTime(c *gin.Context) {
	var req model.CreateOffTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	offTime, err := h.availabilityService.CreateOffTime(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    offTime,
	})
}

func (h *AvailabilityHandler) GetExpertOffTimes(c *gin.Context) {
	expertID, err := strconv.Atoi(c.Param("expert_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expert ID"})
		return
	}

	offTimes, err := h.availabilityService.GetExpertOffTimes(expertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    offTimes,
	})
}

func (h *AvailabilityHandler) DeleteOffTime(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid off time ID"})
		return
	}

	if err := h.availabilityService.DeleteOffTime(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Off time deleted successfully",
	})
}
