package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"booking-service/internal/model"
	"booking-service/internal/service"
	"shared/pkg/logger"
	"shared/pkg/utils"
)

type BookingHandler struct {
	bookingService  service.BookingServiceInterface
	conflictChecker service.ConflictCheckerInterface
	logger          logger.Logger
}

func NewBookingHandler(
	bookingService service.BookingServiceInterface,
	conflictChecker service.ConflictCheckerInterface,
	logger logger.Logger,
) *BookingHandler {
	return &BookingHandler{
		bookingService:  bookingService,
		conflictChecker: conflictChecker,
		logger:          logger,
	}
}

// CreateBooking creates a new booking
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
		return
	}

	var req model.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format"))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	// Check for conflicts
	hasConflict, err := h.conflictChecker.CheckBookingConflict(
		req.ExpertID,
		userID.(uint),
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		h.logger.Error("Failed to check booking conflict", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Internal server error"))
		return
	}

	if hasConflict {
		c.JSON(http.StatusConflict, utils.ErrorResponse("Time slot is already booked"))
		return
	}

	// Create booking
	booking, err := h.bookingService.CreateBooking(userID.(uint), &req)
	if err != nil {
		h.logger.Error("Failed to create booking", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create booking"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Booking created successfully", booking))
}

// GetBooking retrieves a booking by ID
func (h *BookingHandler) GetBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	booking, err := h.bookingService.GetBookingByID(uint(bookingID))
	if err != nil {
		h.logger.Error("Failed to get booking", err)
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		return
	}

	// Check authorization (user can only see their own bookings, experts can see bookings with them)
	if userRole != "admin" && 
		 booking.UserID != userID.(uint) && 
		 booking.ExpertID != userID.(uint) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking retrieved successfully", booking))
}

// UpdateBooking updates a booking
func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	var req model.UpdateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format"))
		return
	}

	// Get existing booking to check authorization
	existingBooking, err := h.bookingService.GetBookingByID(uint(bookingID))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		return
	}

	// Check authorization
	if userRole != "admin" && existingBooking.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		return
	}

	// Check if booking can be updated (not within 1 hour)
	if time.Now().Add(time.Hour).After(existingBooking.StartTime) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Cannot update booking within 1 hour of start time"))
		return
	}

	// Update booking
	booking, err := h.bookingService.UpdateBooking(uint(bookingID), &req)
	if err != nil {
		h.logger.Error("Failed to update booking", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update booking"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking updated successfully", booking))
}

// CancelBooking cancels a booking
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	// Get existing booking to check authorization
	existingBooking, err := h.bookingService.GetBookingByID(uint(bookingID))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		return
	}

	// Check authorization
	if userRole != "admin" && 
		 existingBooking.UserID != userID.(uint) &&
		 existingBooking.ExpertID != userID.(uint) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		return
	}

	// Check if booking can be cancelled (not within 1 hour)
	if time.Now().Add(time.Hour).After(existingBooking.StartTime) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Cannot cancel booking within 1 hour of start time"))
		return
	}

	// Cancel booking
	err = h.bookingService.CancelBooking(uint(bookingID), userID.(uint))
	if err != nil {
		h.logger.Error("Failed to cancel booking", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to cancel booking"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking cancelled successfully", nil))
}

// GetUserBookings retrieves bookings for a user
func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	bookings, total, err := h.bookingService.GetUserBookings(userID.(uint), page, limit, status)
	if err != nil {
		h.logger.Error("Failed to get user bookings", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve bookings"))
		return
	}

	response := map[string]interface{}{
		"bookings": bookings,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Bookings retrieved successfully", response))
}

// GetExpertBookings retrieves bookings for an expert
func (h *BookingHandler) GetExpertBookings(c *gin.Context) {
	expertID, _ := c.Get("user_id")
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	date := c.Query("date")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	bookings, total, err := h.bookingService.GetExpertBookings(expertID.(uint), page, limit, status, date)
	if err != nil {
		h.logger.Error("Failed to get expert bookings", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve bookings"))
		return
	}

	response := map[string]interface{}{
		"bookings": bookings,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Expert bookings retrieved successfully", response))
}