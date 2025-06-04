package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"booking-service/internal/model"
	"booking-service/internal/service"
	"shared/pkg/logger"
	"shared/pkg/utils"
)

type StatusHandler struct {
	statusService service.StatusServiceInterface
	logger        logger.Logger
}

func NewStatusHandler(statusService service.StatusServiceInterface, logger logger.Logger) *StatusHandler {
	return &StatusHandler{
		statusService: statusService,
		logger:        logger,
	}
}

// UpdateBookingStatus updates the status of a booking
func (h *StatusHandler) UpdateBookingStatus(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	var req model.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request format"))
		return
	}

	// Validate status
	if !model.IsValidBookingStatus(req.Status) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking status"))
		return
	}

	// Update status
	err = h.statusService.UpdateBookingStatus(
		uint(bookingID),
		req.Status,
		userID.(uint),
		userRole.(string),
		req.Note,
	)
	if err != nil {
		switch err.Error() {
		case "booking not found":
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		case "access denied":
			c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		case "invalid status transition":
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid status transition"))
		default:
			h.logger.Error("Failed to update booking status", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update status"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking status updated successfully", nil))
}

// GetBookingStatus retrieves the current status of a booking
func (h *StatusHandler) GetBookingStatus(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	status, err := h.statusService.GetBookingStatus(uint(bookingID))
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		} else {
			h.logger.Error("Failed to get booking status", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve status"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking status retrieved successfully", map[string]interface{}{
		"status": status,
	}))
}

// GetStatusHistory retrieves the status history of a booking
func (h *StatusHandler) GetStatusHistory(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	history, err := h.statusService.GetStatusHistory(uint(bookingID), userID.(uint), userRole.(string))
	if err != nil {
		switch err.Error() {
		case "booking not found":
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		case "access denied":
			c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		default:
			h.logger.Error("Failed to get status history", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve status history"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Status history retrieved successfully", history))
}

// ConfirmBooking confirms a booking (expert only)
func (h *StatusHandler) ConfirmBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	var req model.ConfirmBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Note is optional, so we can proceed without it
		req.Note = ""
	}

	err = h.statusService.UpdateBookingStatus(
		uint(bookingID),
		model.BookingStatusConfirmed,
		userID.(uint),
		userRole.(string),
		req.Note,
	)
	if err != nil {
		switch err.Error() {
		case "booking not found":
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		case "access denied":
			c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		case "invalid status transition":
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Cannot confirm this booking"))
		default:
			h.logger.Error("Failed to confirm booking", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to confirm booking"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking confirmed successfully", nil))
}

// RejectBooking rejects a booking (expert only)
func (h *StatusHandler) RejectBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	var req model.RejectBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Rejection reason is required"))
		return
	}

	if req.Reason == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Rejection reason is required"))
		return
	}

	err = h.statusService.UpdateBookingStatus(
		uint(bookingID),
		model.BookingStatusRejected,
		userID.(uint),
		userRole.(string),
		req.Reason,
	)
	if err != nil {
		switch err.Error() {
		case "booking not found":
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		case "access denied":
			c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		case "invalid status transition":
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Cannot reject this booking"))
		default:
			h.logger.Error("Failed to reject booking", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to reject booking"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking rejected successfully", nil))
}

// CompleteBooking marks a booking as completed (expert only)
func (h *StatusHandler) CompleteBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid booking ID"))
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	var req model.CompleteBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Summary is optional
		req.Summary = ""
	}

	err = h.statusService.UpdateBookingStatus(
		uint(bookingID),
		model.BookingStatusCompleted,
		userID.(uint),
		userRole.(string),
		req.Summary,
	)
	if err != nil {
		switch err.Error() {
		case "booking not found":
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Booking not found"))
		case "access denied":
			c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied"))
		case "invalid status transition":
			c.JSON(http.StatusBadRequest, utils.ErrorResponse("Cannot complete this booking"))
		default:
			h.logger.Error("Failed to complete booking", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to complete booking"))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Booking completed successfully", nil))
}