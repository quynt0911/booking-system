package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"services/booking-service/internal/model"
	"services/booking-service/internal/service"
	"services/booking-service/pkg/logger"
	"services/booking-service/pkg/utils"
)

type BookingHandler struct {
	bookingService  service.BookingServiceInterface
	conflictChecker service.ConflictCheckerInterface
	logger          logger.LoggerInterface
}

func NewBookingHandler(
	bookingService service.BookingServiceInterface,
	conflictChecker service.ConflictCheckerInterface,
	logger logger.LoggerInterface,
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

	// Calculate end time from scheduled time and duration
	endTime := req.ScheduledTime.Add(time.Duration(req.DurationMinutes) * time.Minute)

	// Check for conflicts
	hasConflict, err := h.conflictChecker.CheckBookingConflict(
		req.ExpertID,
		userID.(uuid.UUID),
		req.ScheduledTime,
		endTime,
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
	booking, err := h.bookingService.CreateBooking(userID.(uuid.UUID), &req)
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
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	booking, err := h.bookingService.GetBookingByID(bookingID)
	if err != nil {
		h.logger.Error("Failed to get booking", err)
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		return
	}

	// Check authorization (user can only see their own bookings, experts can see bookings with them)
	if userRole != "admin" &&
		booking.UserID != userID.(uuid.UUID) &&
		booking.ExpertID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking retrieved successfully", booking))
}

// UpdateBooking updates a booking
func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
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
	existingBooking, err := h.bookingService.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		return
	}

	// Check authorization
	if userRole != "admin" && existingBooking.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		return
	}

	// Check if booking can be updated (not within 1 hour)
	if time.Now().Add(time.Hour).After(existingBooking.ScheduledTime) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Cannot update booking within 1 hour of start time"))
		return
	}

	// Update booking
	booking, err := h.bookingService.UpdateBooking(bookingID, &req)
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
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	// Get existing booking to check authorization
	existingBooking, err := h.bookingService.GetBookingByID(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		return
	}

	// Check authorization
	if userRole != "admin" &&
		existingBooking.UserID != userID.(uuid.UUID) &&
		existingBooking.ExpertID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		return
	}

	// Check if booking can be cancelled (not within 1 hour)
	if time.Now().Add(time.Hour).After(existingBooking.ScheduledTime) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Cannot cancel booking within 1 hour of start time"))
		return
	}

	var req model.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format"))
		return
	}

	// Cancel booking
	err = h.bookingService.CancelBooking(bookingID, userID.(uuid.UUID), &req)
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

	var req model.GetBookingsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid query parameters"))
		return
	}

	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	bookings, total, err := h.bookingService.GetUserBookings(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to get user bookings", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve bookings"))
		return
	}

	response := model.BookingListResponse{
		Bookings:   bookings,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Bookings retrieved successfully", response))
}

// GetExpertBookings retrieves bookings for an expert
func (h *BookingHandler) GetExpertBookings(c *gin.Context) {
	expertID, _ := c.Get("user_id")

	var req model.GetBookingsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid query parameters"))
		return
	}

	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	bookings, total, err := h.bookingService.GetExpertBookings(expertID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to get expert bookings", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve bookings"))
		return
	}

	response := model.BookingListResponse{
		Bookings:   bookings,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Bookings retrieved successfully", response))
}
