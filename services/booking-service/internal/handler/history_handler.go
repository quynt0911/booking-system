package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"services/booking-service/internal/model"
	"services/booking-service/internal/service"
	"services/booking-service/pkg/logger"
	"services/booking-service/pkg/utils"
)

type HistoryHandler struct {
	bookingService service.BookingServiceInterface
	logger         logger.LoggerInterface
}

func NewHistoryHandler(bookingService service.BookingServiceInterface, logger logger.LoggerInterface) *HistoryHandler {
	return &HistoryHandler{
		bookingService: bookingService,
		logger:         logger,
	}
}

// GetBookingHistory retrieves booking history for a user
func (h *HistoryHandler) GetBookingHistory(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("user_id not found in context")
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		h.logger.Error("user_id is not of type uuid.UUID")
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Internal server error"))
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	statusStr := c.Query("status")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse date filters
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set to end of day
			endOfDay := parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &endOfDay
		}
	}

	// Convert status string to BookingStatus
	var status *model.BookingStatus
	if statusStr != "" {
		bookingStatus := model.BookingStatus(statusStr)
		status = &bookingStatus
	}

	// Create request object
	req := &model.GetHistoryRequest{
		Page:      page,
		Limit:     limit,
		NewStatus: status,
		StartDate: startDate,
		EndDate:   endDate,
	}

	bookings, total, err := h.bookingService.GetBookingHistory(userID, req)
	if err != nil {
		h.logger.Error("Failed to get booking history", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve booking history"))
		return
	}

	response := map[string]interface{}{
		"bookings": bookings,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
		"filters": map[string]interface{}{
			"status":     statusStr,
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking history retrieved successfully", response))
}

// GetExpertHistory retrieves booking history for an expert
func (h *HistoryHandler) GetExpertHistory(c *gin.Context) {
	expertIDVal, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("user_id not found in context")
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
		return
	}
	expertID, ok := expertIDVal.(uuid.UUID)
	if !ok {
		h.logger.Error("user_id is not of type uuid.UUID")
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Internal server error"))
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	statusStr := c.Query("status")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse date filters
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endOfDay := parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &endOfDay
		}
	}

	// Convert status string to BookingStatus
	var status *model.BookingStatus
	if statusStr != "" {
		bookingStatus := model.BookingStatus(statusStr)
		status = &bookingStatus
	}

	// Create request object
	req := &model.GetHistoryRequest{
		Page:      page,
		Limit:     limit,
		NewStatus: status,
		StartDate: startDate,
		EndDate:   endDate,
	}

	bookings, total, err := h.bookingService.GetExpertHistory(expertID, req)
	if err != nil {
		h.logger.Error("Failed to get expert history", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve expert history"))
		return
	}

	response := map[string]interface{}{
		"bookings": bookings,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
		"filters": map[string]interface{}{
			"status":     statusStr,
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Expert history retrieved successfully", response))
}

// GetUpcomingBookings retrieves upcoming bookings for a user or expert
func (h *HistoryHandler) GetUpcomingBookings(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("user_id not found in context")
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		h.logger.Error("user_id is not of type uuid.UUID")
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Internal server error"))
		return
	}

	userRoleVal, exists := c.Get("user_role")
	if !exists {
		h.logger.Error("user_role not found in context")
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
		return
	}
	userRole, ok := userRoleVal.(string)
	if !ok {
		h.logger.Error("user_role is not of type string")
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Internal server error"))
		return
	}

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7")) // next 7 days by default

	if limit < 1 || limit > 50 {
		limit = 10
	}
	if days < 1 || days > 30 {
		days = 7
	}

	var bookings interface{}
	var err error

	if userRole == "expert" {
		bookings, err = h.bookingService.GetUpcomingExpertBookings(userID, limit, days)
	} else {
		bookings, err = h.bookingService.GetUpcomingUserBookings(userID, limit, days)
	}

	if err != nil {
		h.logger.Error("Failed to get upcoming bookings", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve upcoming bookings"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Upcoming bookings retrieved successfully", bookings))
}

// GetPastBookings retrieves past bookings for a user or expert
func (h *HistoryHandler) GetPastBookings(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("user_id not found in context")
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		h.logger.Error("user_id is not of type uuid.UUID")
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Internal server error"))
		return
	}

	userRoleVal, exists := c.Get("user_role")
	if !exists {
		h.logger.Error("user_role not found in context")
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
		return
	}
	userRole, ok := userRoleVal.(string)
	if !ok {
		h.logger.Error("user_role is not of type string")
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Internal server error"))
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var bookings interface{}
	var total int64
	var err error

	if userRole == "expert" {
		bookings, total, err = h.bookingService.GetPastExpertBookings(userID, page, limit)
	} else {
		bookings, total, err = h.bookingService.GetPastUserBookings(userID, page, limit)
	}

	if err != nil {
		h.logger.Error("Failed to get past bookings", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve past bookings"))
		return
	}

	response := map[string]interface{}{
		"bookings": bookings,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Past bookings retrieved successfully", response))
}
