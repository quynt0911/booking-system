package routes

import (
	"expert-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	expertHandler *handler.ExpertHandler,
	scheduleHandler *handler.ScheduleHandler,
	availabilityHandler *handler.AvailabilityHandler,
) {
	// Expert routes
	experts := router.Group("/api/experts")
	{
		experts.POST("", expertHandler.CreateExpert)
		experts.GET("", expertHandler.GetExperts)
		experts.GET("/:id", expertHandler.GetExpert)
		experts.PUT("/:id", expertHandler.UpdateExpert)
		experts.DELETE("/:id", expertHandler.DeleteExpert)
		experts.GET("/expertise", expertHandler.GetExpertsByExpertise)
	}

	// Schedule routes
	schedules := router.Group("/api/schedules")
	{
		schedules.POST("", scheduleHandler.CreateSchedule)
		schedules.GET("", scheduleHandler.GetSchedules)
		schedules.GET("/:id", scheduleHandler.GetScheduleByID)
		schedules.PUT("/:id", scheduleHandler.UpdateSchedule)
		schedules.DELETE("/:id", scheduleHandler.DeleteSchedule)
		schedules.GET("/upcoming", scheduleHandler.GetUpcomingSchedules)
	}

	// Availability routes
	availability := router.Group("/api/v1/availability")
	{
		availability.POST("", availabilityHandler.CreateAvailability)
		availability.GET("", availabilityHandler.GetAvailabilities)
		availability.GET("/:id", availabilityHandler.GetAvailabilityByID)
		availability.PUT("/:id", availabilityHandler.UpdateAvailability)
		availability.DELETE("/:id", availabilityHandler.DeleteAvailability)
		availability.POST("/:id/book", availabilityHandler.BookAvailability)
		availability.POST("/recurring", availabilityHandler.CreateRecurringAvailability)
		availability.POST("/check", availabilityHandler.CheckAvailability)
		availability.POST("/off-time", availabilityHandler.CreateOffTime)
		availability.GET("/off-time/:expert_id", availabilityHandler.GetExpertOffTimes)
		availability.DELETE("/off-time/:id", availabilityHandler.DeleteOffTime)
	}
}
