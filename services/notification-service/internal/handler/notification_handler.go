package handler

import (
	"net/http"
	"notification-service/internal/service"

	"github.com/gin-gonic/gin"
)

func GetRecentNotifications(c *gin.Context) {
	notis := service.FetchRecentNotifications()
	c.JSON(http.StatusOK, notis)
}
