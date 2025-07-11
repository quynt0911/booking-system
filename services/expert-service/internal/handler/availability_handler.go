// services/expert-service/internal/handler/availability_handler.go
package handler

import (
	"encoding/json"
	"expert-service/internal/model"
	"expert-service/internal/service"
	"expert-service/internal/utils" // Thêm import này
	"log"
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

// GetAvailabilities godoc
// @Summary Get filtered availabilities
// @Description Get availability slots with optional filters
// @Tags availability
// @Accept json
// @Produce json
// @Param expert_id query string true "Expert ID"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {array} model.AvailabilitySlot
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

	// Validate: end_date must be >= start_date
	if endDate.Before(startDate) {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "end_date must be greater than or equal to start_date"})
		return
	}

	slots, err := h.availabilityService.GetAvailabilities(expertID, startDate, endDate, nil, "")
	if err != nil {
		log.Printf("Error getting availabilities: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, slots)
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

	// Sau khi tạo off-time, sinh lại availability cho 14 ngày tới
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 14)

	// Tạo service token cho internal call
	serviceToken, err := utils.GenerateServiceToken("expert-service")
	if err != nil {
		log.Printf("Error generating service token: %v", err)
		// Không return error ở đây vì off-time đã được tạo thành công
	} else {
		slots, err := h.availabilityService.GetAvailabilities(req.ExpertID, startDate, endDate, nil, serviceToken)
		if err == nil {
			for _, slot := range slots {
				key := "availability:" + req.ExpertID + ":" + slot.Date
				data, _ := json.Marshal(slot)
				h.availabilityService.(interface {
					SetAvailability(key string, value []byte) error
				}).SetAvailability(key, data)
			}
		}
	}

	c.JSON(http.StatusCreated, offTime)
}

// Các method khác giữ nguyên không đổi...
// CreateAvailability, GetAvailabilityByID, UpdateAvailability, DeleteAvailability,
// BookAvailability, CreateRecurringAvailability, CheckAvailability,
// GetExpertOffTimes, DeleteOffTime, RegisterRoutes

// CreateAvailability godoc
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

// GetAvailabilityByID godoc
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
func (h *AvailabilityHandler) DeleteAvailability(c *gin.Context) {
	id := c.Param("id")

	if err := h.availabilityService.DeleteAvailability(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// BookAvailability godoc
func (h *AvailabilityHandler) BookAvailability(c *gin.Context) {
	id := c.Param("id")

	if err := h.availabilityService.BookAvailability(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// CreateRecurringAvailability godoc
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

// GetExpertOffTimes godoc
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
