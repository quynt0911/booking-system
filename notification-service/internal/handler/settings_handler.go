package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateNotificationSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}
