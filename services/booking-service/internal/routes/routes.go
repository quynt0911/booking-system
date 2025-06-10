package routes

import (
	"services/booking-service/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRoutes thiết lập các route cho booking service
func SetupRoutes(router *gin.Engine, bookingHandler *handler.BookingHandler, statusHandler *handler.StatusHandler, historyHandler *handler.HistoryHandler) {
	// Booking routes
	router.POST("/", bookingHandler.CreateBooking)
	router.GET("/:id", bookingHandler.GetBooking)
	router.PUT("/:id", bookingHandler.UpdateBooking)
	router.DELETE("/:id", bookingHandler.CancelBooking)
	router.GET("/user", bookingHandler.GetUserBookings)     // Adjusted path
	router.GET("/expert", bookingHandler.GetExpertBookings) // Adjusted path

	// Status routes
	router.PUT("/:id/status", statusHandler.UpdateBookingStatus)
	router.GET("/:id/status", statusHandler.GetBookingStatus)
	router.GET("/:id/status/history", statusHandler.GetStatusHistory)

	// History routes
	router.GET("/history/user", historyHandler.GetBookingHistory)
	router.GET("/history/expert", historyHandler.GetExpertHistory)
}
