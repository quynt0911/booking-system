package handler

import (
	"expert-service/internal/model"
	"expert-service/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ScheduleHandler handles schedule-related requests
type ScheduleHandler struct {
	scheduleService service.ScheduleService
}

// NewScheduleHandler creates a new schedule handler
func NewScheduleHandler(scheduleService service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

// Schedule represents a scheduled appointment
type Schedule struct {
	ID             string    `json:"id"`
	ExpertID       string    `json:"expert_id"`
	UserID         string    `json:"user_id"`
	AvailabilityID string    `json:"availability_id"`
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
// THIS STRUCT IS REMOVED, USING model.CreateScheduleRequest INSTEAD

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
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req model.CreateScheduleRequest // Using the model's struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	schedule, err := h.scheduleService.CreateSchedule(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    schedule,
		"message": "Schedule created successfully",
	})
}

// GetSchedules retrieves all schedules with optional filters
func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
	var req model.GetSchedulesRequest
	req.ExpertID = c.Query("expert_id")
	req.UserID = c.Query("user_id")
	req.Status = c.Query("status")
	req.Date = c.Query("date")

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	schedules, err := h.scheduleService.GetSchedules(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve schedules", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
		"filters": gin.H{
			"expert_id": req.ExpertID,
			"user_id":   req.UserID,
			"status":    req.Status,
			"date":      req.Date,
			"limit":     req.Limit,
			"offset":    req.Offset,
		},
	})
}

// GetScheduleByID retrieves schedule by ID
func (h *ScheduleHandler) GetScheduleByID(c *gin.Context) {
	id := c.Param("id")                       // ID is already a string (UUID)
	if _, err := uuid.Parse(id); err != nil { // Validate if it's a valid UUID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID format"})
		return
	}

	schedule, err := h.scheduleService.GetScheduleByID(uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve schedule", "details": err.Error()})
		return
	}

	if schedule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedule,
	})
}

// UpdateSchedule updates existing schedule
func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id := c.Param("id")                       // ID is already a string (UUID)
	if _, err := uuid.Parse(id); err != nil { // Validate if it's a valid UUID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID format"})
		return
	}

	var req model.UpdateScheduleRequest // Use the model's struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if req.Status != nil && !isValidStatus(*req.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Use: pending, confirmed, cancelled, completed"})
		return
	}

	if err := h.scheduleService.UpdateSchedule(uuid.MustParse(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Schedule updated successfully",
	})
}

// CancelSchedule cancels a schedule
func (h *ScheduleHandler) CancelSchedule(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID format"})
		return
	}

	if err := h.scheduleService.CancelSchedule(uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel schedule", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Schedule cancelled successfully",
	})
}

// ConfirmSchedule confirms a pending schedule
func (h *ScheduleHandler) ConfirmSchedule(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID format"})
		return
	}

	if err := h.scheduleService.ConfirmSchedule(uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm schedule", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Schedule confirmed successfully",
	})
}

// CompleteSchedule marks a schedule as completed
func (h *ScheduleHandler) CompleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID format"})
		return
	}

	if err := h.scheduleService.CompleteSchedule(uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete schedule", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Schedule marked as completed",
	})
}

// GetExpertSchedules retrieves all schedules for a specific expert
func (h *ScheduleHandler) GetExpertSchedules(c *gin.Context) {
	expertID := c.Param("expert_id")
	if _, err := uuid.Parse(expertID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expert ID format"})
		return
	}

	var req model.GetSchedulesRequest
	req.ExpertID = expertID
	req.Status = c.Query("status")
	req.Date = c.Query("date")

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	schedules, err := h.scheduleService.GetSchedules(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expert schedules", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      schedules,
		"expert_id": expertID,
		"filters": gin.H{
			"status": req.Status,
			"date":   req.Date,
			"limit":  req.Limit,
			"offset": req.Offset,
		},
	})
}

// GetUserSchedules retrieves all schedules for a specific user
func (h *ScheduleHandler) GetUserSchedules(c *gin.Context) {
	userID := c.Param("user_id")
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req model.GetSchedulesRequest
	req.UserID = userID
	req.Status = c.Query("status")
	req.Date = c.Query("date")

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	schedules, err := h.scheduleService.GetSchedules(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user schedules", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
		"user_id": userID,
		"filters": gin.H{
			"status": req.Status,
			"date":   req.Date,
			"limit":  req.Limit,
			"offset": req.Offset,
		},
	})
}

// GetTodaySchedules retrieves schedules for today
func (h *ScheduleHandler) GetTodaySchedules(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	var req model.GetSchedulesRequest
	req.Date = today
	req.ExpertID = c.Query("expert_id") // Optional filter
	req.Status = c.Query("status")      // Optional filter

	schedules, err := h.scheduleService.GetSchedules(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve today's schedules", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
		"date":    today,
		"filters": gin.H{
			"expert_id": req.ExpertID,
			"status":    req.Status,
		},
	})
}

// DeleteSchedule deletes a schedule
func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")              // ID is already a string (UUID)
	offTimeID, err := uuid.Parse(id) // Use offTimeID instead of id for consistency with the Availability handler
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID format"})
		return
	}

	if err := h.scheduleService.DeleteSchedule(offTimeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Schedule deleted successfully",
	})
}

// GetUpcomingSchedules retrieves upcoming schedules for a specific expert
func (h *ScheduleHandler) GetUpcomingSchedules(c *gin.Context) {
	expertID := c.Query("expert_id")
	daysStr := c.Query("days")

	days := 7 // Default to 7 days
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")

	var req model.GetSchedulesRequest
	req.ExpertID = expertID
	// For upcoming schedules, we typically look for pending or confirmed
	req.Status = "pending" // Or fetch both pending and confirmed
	// req.Status = "confirmed"
	// We need a way to filter by date range in GetSchedules, for now, we'll just pass expertID.
	// A more robust solution would involve adding startDate/endDate to GetSchedulesRequest and repository.

	schedules, err := h.scheduleService.GetSchedules(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve upcoming schedules", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      schedules,
		"expert_id": expertID,
		"days":      days,
		"filters": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
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
