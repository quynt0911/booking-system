package handler

import (
	"net/http"

	"github.com/your-org/booking-system/services/notification-service/internal/service"

	"github.com/gin-gonic/gin"
)

// NotificationHandler struct
type NotificationHandler struct {
	EmailService *service.EmailService
}

func NewNotificationHandler(emailService *service.EmailService) *NotificationHandler {
	return &NotificationHandler{EmailService: emailService}
}

func (h *NotificationHandler) SendEmailHandler(c *gin.Context) {
	var req struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := h.EmailService.Send(req.To, req.Subject, req.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func GetRecentNotifications(c *gin.Context) {
	notis := service.FetchRecentNotifications()
	c.JSON(http.StatusOK, notis)
}
