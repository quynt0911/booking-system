package routes

import (
	"expert-service/internal/handler"
	"expert-service/internal/middleware"
	"expert-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	expertService service.ExpertService,
	expertScheduleService service.ExpertScheduleService,
	availabilityService service.ExpertAvailabilityService,
) {
	// Initialize handlers
	expertHandler := handler.NewExpertHandler(expertService)
	expertScheduleHandler := handler.NewExpertScheduleHandler(expertScheduleService, availabilityService)
	availabilityHandler := handler.NewAvailabilityHandler(availabilityService)
	offTimeHandler := handler.NewOffTimeHandler(availabilityService, expertService)

	// Expert routes
	experts := router.Group("/experts", middleware.JWTAuthMiddleware())
	{
		experts.POST("", expertHandler.CreateExpert)
		experts.GET("", expertHandler.GetExperts)
		experts.GET("/:id", expertHandler.GetExpert)
		experts.PUT("/:id", expertHandler.UpdateExpert)
		experts.DELETE("/:id", expertHandler.DeleteExpert)
		experts.GET("/expertise", expertHandler.GetExpertsByExpertise)
	}

	// Availability routes
	availability := router.Group("/availability")
	{
		availability.GET("", availabilityHandler.GetAvailabilities)
		availability.GET("/:id", availabilityHandler.GetAvailabilityByID)
		availability.POST("/check", availabilityHandler.CheckAvailability)
	}

	// OffTime routes (tách riêng)
	offTime := router.Group("/off-time")
	offTime.Use(middleware.JWTAuthMiddleware())
	{
		offTime.POST("", offTimeHandler.CreateOffTime)
		offTime.GET(":expert_id", offTimeHandler.GetExpertOffTimes)
		offTime.GET("/detail/:id", offTimeHandler.GetOffTimeByID)
		offTime.DELETE(":id", offTimeHandler.DeleteOffTime)
	}

	// Expert Schedule routes
	expertSchedules := router.Group("/expert-schedules")
	{
		expertSchedules.POST("", expertScheduleHandler.CreateSchedule)
		expertSchedules.GET("", expertScheduleHandler.GetSchedulesByExpertID)
		expertSchedules.PUT("/:id", expertScheduleHandler.UpdateSchedule)
		expertSchedules.DELETE("/:id", expertScheduleHandler.DeleteSchedule)
	}
}
