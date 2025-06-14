package handler

import (
	"encoding/json"
	"expert-service/internal/model"
	"expert-service/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Message: message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// Availability represents expert availability
type Availability struct {
	ID        string    `json:"id"`
	ExpertID  string    `json:"expert_id"`
	Date      string    `json:"date"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	IsBooked  bool      `json:"is_booked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateAvailabilityRequest represents request body for creating availability
type CreateAvailabilityRequest struct {
	ExpertID  string `json:"expert_id" validate:"required"`
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

// CreateAvailability godoc
// @Summary Create a new availability slot
// @Description Create a new availability slot for an expert
// @Tags availability
// @Accept json
// @Produce json
// @Param availability body model.CreateAvailabilityRequest true "Availability details"
// @Success 201 {object} model.Availability
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability [post]
func (h *AvailabilityHandler) CreateAvailability(c *gin.Context) {
	var req model.CreateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request body"})
		return
	}

	availability, err := h.availabilityService.CreateAvailability(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, availability)
}

// GetAvailabilities godoc
// @Summary Get filtered availabilities
// @Description Get availability slots with optional filters
// @Tags availability
// @Accept json
// @Produce json
// @Param expert_id query string true "Expert ID"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param is_booked query bool false "Filter by booking status"
// @Success 200 {array} model.Availability
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability [get]
func (h *AvailabilityHandler) GetAvailabilities(c *gin.Context) {
	expertID := c.Query("expert_id")
	if expertID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Expert ID is required"})
		return
	}

	startDateStr := c.Query("start_date")
	if startDateStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Start date is required"})
		return
	}
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid start date format"})
		return
	}

	endDateStr := c.Query("end_date")
	if endDateStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "End date is required"})
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid end date format"})
		return
	}

	var isBooked *bool
	if isBookedStr := c.Query("is_booked"); isBookedStr != "" {
		booked := isBookedStr == "true"
		isBooked = &booked
	}

	availabilities, err := h.availabilityService.GetAvailabilities(expertID, startDate, endDate, isBooked)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, availabilities)
}

// GetAvailabilityByID godoc
// @Summary Get availability by ID
// @Description Get a specific availability slot by its ID
// @Tags availability
// @Accept json
// @Produce json
// @Param id path string true "Availability ID"
// @Success 200 {object} model.Availability
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability/{id} [get]
func (h *AvailabilityHandler) GetAvailabilityByID(c *gin.Context) {
	id := c.Param("id")

	availability, err := h.availabilityService.GetAvailabilityByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	if availability == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Availability not found"})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// UpdateAvailability godoc
// @Summary Update availability
// @Description Update an existing availability slot
// @Tags availability
// @Accept json
// @Produce json
// @Param id path string true "Availability ID"
// @Param availability body model.UpdateAvailabilityRequest true "Updated availability details"
// @Success 200 {object} model.Availability
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability/{id} [put]
func (h *AvailabilityHandler) UpdateAvailability(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request body"})
		return
	}

	availability, err := h.availabilityService.UpdateAvailability(id, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	if availability == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Message: "Availability not found"})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// DeleteAvailability godoc
// @Summary Delete availability
// @Description Delete an availability slot
// @Tags availability
// @Accept json
// @Produce json
// @Param id path string true "Availability ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability/{id} [delete]
func (h *AvailabilityHandler) DeleteAvailability(c *gin.Context) {
	id := c.Param("id")

	if err := h.availabilityService.DeleteAvailability(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// BookAvailability godoc
// @Summary Book availability
// @Description Book an availability slot
// @Tags availability
// @Accept json
// @Produce json
// @Param id path string true "Availability ID"
// @Success 200 {object} model.Availability
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability/{id}/book [post]
func (h *AvailabilityHandler) BookAvailability(c *gin.Context) {
	id := c.Param("id")

	if err := h.availabilityService.BookAvailability(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// CreateRecurringAvailability godoc
// @Summary Create recurring availability
// @Description Create multiple availability slots for recurring schedules
// @Tags availability
// @Accept json
// @Produce json
// @Param availability body model.CreateRecurringAvailabilityRequest true "Recurring availability details"
// @Success 201 {array} model.Availability
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability/recurring [post]
func (h *AvailabilityHandler) CreateRecurringAvailability(c *gin.Context) {
	var req model.CreateRecurringAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request body"})
		return
	}

	availabilities, err := h.availabilityService.CreateRecurringAvailability(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, availabilities)
}

// CheckAvailability godoc
// @Summary Check availability for a time slot
// @Description Check if an expert is available at a specific time slot
// @Tags availability
// @Accept json
// @Produce json
// @Param availability body model.CheckAvailabilityRequest true "Availability check details"
// @Success 200 {object} bool "true if available, false otherwise"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability/check [post]
func (h *AvailabilityHandler) CheckAvailability(c *gin.Context) {
	var req model.CheckAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request body"})
		return
	}

	isAvailable, err := h.availabilityService.CheckAvailability(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, isAvailable)
}

// CreateOffTime godoc
// @Summary Create off-time for an expert
// @Description Create a period when an expert is unavailable
// @Tags availability
// @Accept json
// @Produce json
// @Param off_time body model.CreateOffTimeRequest true "Off-time details"
// @Success 201 {object} model.OffTime
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability/off-time [post]
func (h *AvailabilityHandler) CreateOffTime(c *gin.Context) {
	var req model.CreateOffTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request body"})
		return
	}

	offTime, err := h.availabilityService.CreateOffTime(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, offTime)
}

// GetExpertOffTimes godoc
// @Summary Get off-times for an expert
// @Description Get all off-time periods for a specific expert
// @Tags availability
// @Accept json
// @Produce json
// @Param expert_id path string true "Expert ID"
// @Success 200 {array} model.OffTime
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability/off-time/{expert_id} [get]
func (h *AvailabilityHandler) GetExpertOffTimes(c *gin.Context) {
	expertID := c.Param("expert_id")
	if expertID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Expert ID is required"})
		return
	}

	offTimes, err := h.availabilityService.GetExpertOffTimes(expertID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, offTimes)
}

// DeleteOffTime godoc
// @Summary Delete an off-time entry
// @Description Delete a specific off-time entry by its ID
// @Tags availability
// @Accept json
// @Produce json
// @Param id path string true "Off-time ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/availability/off-time/{id} [delete]
func (h *AvailabilityHandler) DeleteOffTime(c *gin.Context) {
	id := c.Param("id")

	if err := h.availabilityService.DeleteOffTime(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RegisterRoutes registers all availability routes
func (h *AvailabilityHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/availability", h.CreateAvailability)
	router.GET("/availability", h.GetAvailabilities)
	router.GET("/availability/:id", h.GetAvailabilityByID)
	router.PUT("/availability/:id", h.UpdateAvailability)
	router.DELETE("/availability/:id", h.DeleteAvailability)
	router.POST("/availability/:id/book", h.BookAvailability)
	router.POST("/availability/recurring", h.CreateRecurringAvailability)
	router.POST("/availability/check", h.CheckAvailability)
	router.POST("/availability/off-time", h.CreateOffTime)
	router.GET("/availability/off-time/:expert_id", h.GetExpertOffTimes)
	router.DELETE("/availability/off-time/:id", h.DeleteOffTime)
}
