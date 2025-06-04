package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// ScheduleHandler handles schedule-related requests
type ScheduleHandler struct {
	// Add your dependencies here (e.g., service, repository)
}

// NewScheduleHandler creates a new schedule handler
func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{}
}

// Schedule represents a scheduled appointment
type Schedule struct {
	ID             int       `json:"id"`
	ExpertID       int       `json:"expert_id"`
	UserID         int       `json:"user_id"`
	AvailabilityID int       `json:"availability_id"`
	Date           string    `json:"date"`
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	Status         string    `json:"status"` // pending, confirmed, cancelled, completed
	Title          string    `json:"title"`
	Description    string    `json:"description,omitempty"`
	MeetingLink    string    `json:"meeting_link,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ScheduleWithDetails includes expert and user information
type ScheduleWithDetails struct {
	Schedule
	ExpertName string `json:"expert_name"`
	UserName   string `json:"user_name"`
	UserEmail  string `json:"user_email"`
}

// CreateScheduleRequest represents request body for creating schedule
type CreateScheduleRequest struct {
	ExpertID       int    `json:"expert_id" validate:"required"`
	UserID         int    `json:"user_id" validate:"required"`
	AvailabilityID int    `json:"availability_id" validate:"required"`
	Title          string `json:"title" validate:"required"`
	Description    string `json:"description,omitempty"`
	MeetingLink    string `json:"meeting_link,omitempty"`
}

// UpdateScheduleRequest represents request body for updating schedule
type UpdateScheduleRequest struct {
	Status      string `json:"status,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	MeetingLink string `json:"meeting_link,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

// ScheduleStatus constants
const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)

// CreateSchedule creates a new schedule/appointment
func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var req CreateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to:
	// 1. Check if availability exists and is not booked
	// 2. Check if expert exists
	// 3. Check if user exists
	// 4. Create schedule and mark availability as booked

	// Example response
	schedule := Schedule{
		ID:             1,
		ExpertID:       req.ExpertID,
		UserID:         req.UserID,
		AvailabilityID: req.AvailabilityID,
		Date:           "2024-12-10",
		StartTime:      "09:00",
		EndTime:        "10:00",
		Status:         StatusPending,
		Title:          req.Title,
		Description:    req.Description,
		MeetingLink:    req.MeetingLink,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    schedule,
		"message": "Schedule created successfully",
	})
}

// GetSchedules retrieves all schedules with optional filters
func (h *ScheduleHandler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	expertID := r.URL.Query().Get("expert_id")
	userID := r.URL.Query().Get("user_id")
	status := r.URL.Query().Get("status")
	date := r.URL.Query().Get("date")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// TODO: Add business logic to fetch schedules with filters
	// Example response
	schedules := []ScheduleWithDetails{
		{
			Schedule: Schedule{
				ID:             1,
				ExpertID:       1,
				UserID:         1,
				AvailabilityID: 1,
				Date:           "2024-12-10",
				StartTime:      "09:00",
				EndTime:        "10:00",
				Status:         StatusConfirmed,
				Title:          "Consultation Session",
				Description:    "Initial consultation",
				MeetingLink:    "https://meet.google.com/abc-def-ghi",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			ExpertName: "Dr. John Smith",
			UserName:   "Jane Doe",
			UserEmail:  "jane@example.com",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    schedules,
		"filters": map[string]string{
			"expert_id": expertID,
			"user_id":   userID,
			"status":    status,
			"date":      date,
			"limit":     limit,
			"offset":    offset,
		},
	})
}

// GetScheduleByID retrieves schedule by ID
func (h *ScheduleHandler) GetScheduleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to fetch schedule by ID with details
	// Example response
	schedule := ScheduleWithDetails{
		Schedule: Schedule{
			ID:             id,
			ExpertID:       1,
			UserID:         1,
			AvailabilityID: 1,
			Date:           "2024-12-10",
			StartTime:      "09:00",
			EndTime:        "10:00",
			Status:         StatusConfirmed,
			Title:          "Consultation Session",
			Description:    "Initial consultation",
			MeetingLink:    "https://meet.google.com/abc-def-ghi",
			Notes:          "Patient is preparing for surgery",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		ExpertName: "Dr. John Smith",
		UserName:   "Jane Doe",
		UserEmail:  "jane@example.com",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    schedule,
	})
}

// UpdateSchedule updates existing schedule
func (h *ScheduleHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	var req UpdateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate status if provided
	if req.Status != "" && !isValidStatus(req.Status) {
		http.Error(w, "Invalid status. Use: pending, confirmed, cancelled, completed", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to update schedule
	// Handle status changes (e.g., if cancelled, free up availability)

	// Example response
	schedule := Schedule{
		ID:             id,
		ExpertID:       1,
		UserID:         1,
		AvailabilityID: 1,
		Date:           "2024-12-10",
		StartTime:      "09:00",
		EndTime:        "10:00",
		Status:         getStringValueOrDefault(req.Status, StatusConfirmed),
		Title:          getStringValueOrDefault(req.Title, "Consultation Session"),
		Description:    req.Description,
		MeetingLink:    req.MeetingLink,
		Notes:          req.Notes,
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UpdatedAt:      time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    schedule,
		"message": "Schedule updated successfully",
	})
}

// CancelSchedule cancels a schedule
func (h *ScheduleHandler) CancelSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to:
	// 1. Update schedule status to cancelled
	// 2. Free up the availability slot
	// 3. Send notification to both expert and user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Schedule cancelled successfully",
		"id":      id,
	})
}

// ConfirmSchedule confirms a pending schedule
func (h *ScheduleHandler) ConfirmSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to:
	// 1. Update schedule status to confirmed
	// 2. Send confirmation notification
	// 3. Generate meeting link if not provided

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Schedule confirmed successfully",
		"id":      id,
	})
}

// CompleteSchedule marks a schedule as completed
func (h *ScheduleHandler) CompleteSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	// TODO: Add business logic to:
	// 1. Update schedule status to completed
	// 2. Allow expert to add notes/summary
	// 3. Trigger billing/payment if applicable

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Schedule marked as completed",
		"id":      id,
	})
}

// GetExpertSchedules retrieves all schedules for a specific expert
func (h *ScheduleHandler) GetExpertSchedules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	expertID, err := strconv.Atoi(vars["expert_id"])
	if err != nil {
		http.Error(w, "Invalid expert ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	status := r.URL.Query().Get("status")
	date := r.URL.Query().Get("date")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// TODO: Add business logic to fetch expert's schedules
	// Example response
	schedules := []ScheduleWithDetails{
		{
			Schedule: Schedule{
				ID:             1,
				ExpertID:       expertID,
				UserID:         1,
				AvailabilityID: 1,
				Date:           "2024-12-10",
				StartTime:      "09:00",
				EndTime:        "10:00",
				Status:         StatusConfirmed,
				Title:          "Consultation Session",
				Description:    "Initial consultation",
				MeetingLink:    "https://meet.google.com/abc-def-ghi",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			UserName:  "Jane Doe",
			UserEmail: "jane@example.com",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"data":      schedules,
		"expert_id": expertID,
		"filters": map[string]string{
			"status": status,
			"date":   date,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetUserSchedules retrieves all schedules for a specific user
func (h *ScheduleHandler) GetUserSchedules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	status := r.URL.Query().Get("status")
	date := r.URL.Query().Get("date")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// TODO: Add business logic to fetch user's schedules
	// Example response
	schedules := []ScheduleWithDetails{
		{
			Schedule: Schedule{
				ID:             1,
				ExpertID:       1,
				UserID:         userID,
				AvailabilityID: 1,
				Date:           "2024-12-10",
				StartTime:      "09:00",
				EndTime:        "10:00",
				Status:         StatusConfirmed,
				Title:          "Consultation Session",
				Description:    "Initial consultation",
				MeetingLink:    "https://meet.google.com/abc-def-ghi",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			ExpertName: "Dr. John Smith",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    schedules,
		"user_id": userID,
		"filters": map[string]string{
			"status": status,
			"date":   date,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetTodaySchedules retrieves schedules for today
func (h *ScheduleHandler) GetTodaySchedules(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")

	// Parse optional filters
	expertID := r.URL.Query().Get("expert_id")
	status := r.URL.Query().Get("status")

	// TODO: Add business logic to fetch today's schedules
	// Example response
	schedules := []ScheduleWithDetails{
		{
			Schedule: Schedule{
				ID:             1,
				ExpertID:       1,
				UserID:         1,
				AvailabilityID: 1,
				Date:           today,
				StartTime:      "09:00",
				EndTime:        "10:00",
				Status:         StatusConfirmed,
				Title:          "Morning Consultation",
				Description:    "Follow-up session",
				MeetingLink:    "https://meet.google.com/abc-def-ghi",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			ExpertName: "Dr. John Smith",
			UserName:   "Jane Doe",
			UserEmail:  "jane@example.com",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    schedules,
		"date":    today,
		"filters": map[string]string{
			"expert_id": expertID,
			"status":    status,
		},
	})
}

// Helper functions
func isValidStatus(status string) bool {
	validStatuses := []string{StatusPending, StatusConfirmed, StatusCancelled, StatusCompleted}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

func getStringValueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}