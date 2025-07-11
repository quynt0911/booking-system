package handler

import (
	"encoding/json"
	"expert-service/internal/model"
	"expert-service/internal/service"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExpertScheduleHandler struct {
	service             service.ExpertScheduleService
	availabilityService service.ExpertAvailabilityService
}

func NewExpertScheduleHandler(s service.ExpertScheduleService, a service.ExpertAvailabilityService) *ExpertScheduleHandler {
	return &ExpertScheduleHandler{service: s, availabilityService: a}
}

func (h *ExpertScheduleHandler) CreateSchedule(c *gin.Context) {
	var req model.ExpertSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = uuid.New()
	if err := h.service.CreateSchedule(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Sau khi tạo schedule, sinh availability cho 14 ngày tới
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 14)
	token := c.GetHeader("Authorization")
	if token != "" && strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	slots, err := h.availabilityService.GetAvailabilities(req.ExpertID.String(), startDate, endDate, nil, token)
	if err == nil {
		for _, slot := range slots {
			key := "availability:" + req.ExpertID.String() + ":" + slot.Date
			data, _ := json.Marshal(slot)
			if setter, ok := h.availabilityService.(interface {
				SetAvailability(key string, value []byte) error
			}); ok {
				setter.SetAvailability(key, data)
			}
		}
	}

	c.JSON(http.StatusCreated, req)
}

func (h *ExpertScheduleHandler) GetSchedulesByExpertID(c *gin.Context) {
	expertIDStr := c.Query("expert_id")
	expertID, err := uuid.Parse(expertIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expert_id"})
		return
	}
	schedules, err := h.service.GetSchedulesByExpertID(expertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

func (h *ExpertScheduleHandler) UpdateSchedule(c *gin.Context) {
	var req model.ExpertSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateSchedule(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *ExpertScheduleHandler) DeleteSchedule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteSchedule(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
