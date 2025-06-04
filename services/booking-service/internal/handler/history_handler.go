package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"booking-service/internal/service"
	"shared/pkg/logger"
	"shared/pkg/utils"
)

type HistoryHandler struct {
	bookingService service.BookingServiceInterface
	logger         logger.Logger
}

func NewHistoryHandler(bookingService service.BookingServiceInterface, logger logger.Logger) *HistoryHandler {
	return &HistoryHandler{
		bookingService: bookingService,
		logger:         logger,
	}
}

// GetBookingHistory retrieves booking history for a user
func (h *HistoryHandler) GetBookingHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
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

	bookings, total, err := h.bookingService.GetBookingHistory(
		userID.(uint),
		page,
		limit,
		status,
		startDate,
		endDate,
	)
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
			"total_pages": (total + limit - 1) / limit,
		},
		"filters": map[string]interface{}{
			"status":     status,
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking history retrieved successfully", response))
}

// GetExpertHistory retrieves booking history for an expert
func (h *HistoryHandler) GetExpertHistory(c *gin.Context) {
	expertID, _ := c.Get("user_id")
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
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

	bookings, total, err := h.bookingService.GetExpertHistory(
		expertID.(uint),
		page,
		limit,
		status,
		startDate,
		endDate,
	)
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
			"total_pages": (total + limit - 1) / limit,
		},
		"filters": map[string]interface{}{
			"status":     status,
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Expert history retrieved successfully", response))
}

// GetBookingStatistics retrieves booking statistics for a user or expert
func (h *HistoryHandler) GetBookingStatistics(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	
	// Parse query parameters
	period := c.DefaultQuery("period", "month") // week, month, year
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	stats, err := h.bookingService.GetBookingStatistics(
		userID.(uint),
		userRole.(string),
		period,
		year,
		month,
	)
	if err != nil {
		h.logger.Error("Failed to get booking statistics", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve statistics"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Statistics retrieved successfully", stats))
}

// GetUpcomingBookings retrieves upcoming bookings for a user
func (h *HistoryHandler) GetUpcomingBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	
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
		bookings, err = h.bookingService.GetUpcomingExpertBookings(userID.(uint), limit, days)
	} else {
		bookings, err = h.bookingService.GetUpcomingUserBookings(userID.(uint), limit, days)
	}

	if err != nil {
		h.logger.Error("Failed to get upcoming bookings", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve upcoming bookings"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Upcoming bookings retrieved successfully", bookings))
}

// GetPastBookings retrieves past bookings for a user
func (h *HistoryHandler) GetPastBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	
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
	var total int
	var err error

	if userRole == "expert" {
		bookings, total, err = h.bookingService.GetPastExpertBookings(userID.(uint), page, limit)
	} else {
		bookings, total, err = h.bookingService.GetPastUserBookings(userID.(uint), page, limit)
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
			"total_pages": (total + limit - 1) / limit,
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Past bookings retrieved successfully", response))
}