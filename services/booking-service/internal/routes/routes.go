package routes

import (
	"services/booking-service/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRoutes thiết lập các route cho booking service
func SetupRoutes(router *gin.Engine, bookingHandler *handler.BookingHandler, statusHandler *handler.StatusHandler, historyHandler *handler.HistoryHandler) {
	// Booking routes
	router.POST("/CreateBooking", bookingHandler.CreateBooking)
	router.GET("/GetBooking/:id", bookingHandler.GetBooking)
	router.PUT("/UpdateBooking/:id", bookingHandler.UpdateBooking)
	router.DELETE("/DeleteBooking/:id", bookingHandler.CancelBooking)
	router.GET("/GetUserBookings", bookingHandler.GetUserBookings)
	router.GET("/GetExpertBookings", bookingHandler.GetExpertBookings)

	// Status routes
	router.PUT("/UpdateBookingStatus/:id", statusHandler.UpdateBookingStatus)
	router.GET("/GetBookingStatus/:id", statusHandler.GetBookingStatus)
	router.GET("/GetStatusHistory/:id", statusHandler.GetStatusHistory)

	// History routes
	router.GET("/GetBookingHistoryByUser", historyHandler.GetBookingHistory)
	router.GET("/GetBookingHistoryByExpert", historyHandler.GetExpertHistory)
}
