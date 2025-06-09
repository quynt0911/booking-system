package routes

import (
	"github.com/gin-gonic/gin"
	"booking-system/services/booking-service/internal/handler"
)

// SetupRoutes thiết lập các route cho booking service
func SetupRoutes(router *gin.Engine, bookingHandler *handler.BookingHandler, statusHandler *handler.StatusHandler, historyHandler *handler.HistoryHandler) {
	bookingRoutes := router.Group("/bookings")
	{
		bookingRoutes.POST("/", bookingHandler.CreateBooking)
		bookingRoutes.GET("/:id", bookingHandler.GetBooking)
		bookingRoutes.PUT("/:id", bookingHandler.UpdateBooking)
		bookingRoutes.DELETE("/:id", bookingHandler.CancelBooking)
		bookingRoutes.GET("/user/:user_id", bookingHandler.GetUserBookings)
		bookingRoutes.GET("/expert/:expert_id", bookingHandler.GetExpertBookings)
	}

	statusRoutes := router.Group("/status")
	{
		statusRoutes.PUT("/:id", statusHandler.UpdateBookingStatus)
		statusRoutes.GET("/:id", statusHandler.GetBookingStatus)
		statusRoutes.GET("/:id/history", statusHandler.GetStatusHistory)
	}

	historyRoutes := router.Group("/history")
	{
		historyRoutes.GET("/user/:user_id", historyHandler.GetBookingHistory)
		historyRoutes.GET("/expert/:expert_id", historyHandler.GetExpertHistory)
	}
}