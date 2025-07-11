package handler

import (
	"expert-service/internal/model"
	"expert-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OffTimeHandler struct {
	offTimeService service.ExpertAvailabilityService
	expertService  service.ExpertService
}

func NewOffTimeHandler(offTimeService service.ExpertAvailabilityService, expertService service.ExpertService) *OffTimeHandler {
	return &OffTimeHandler{
		offTimeService: offTimeService,
		expertService:  expertService,
	}
}

// Tạo off-time mới
func (h *OffTimeHandler) CreateOffTime(c *gin.Context) {
	var req model.CreateOffTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Lấy user_id từ JWT
	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing user_id in token"})
		return
	}

	// Lấy expert_id từ user_id
	expertID, err := h.expertService.GetExpertIDByUserID(userID.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Gán expert_id vào request
	req.ExpertID = expertID

	offTime, err := h.offTimeService.CreateOffTime(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, offTime)
}

// Lấy danh sách off-time của expert
func (h *OffTimeHandler) GetExpertOffTimes(c *gin.Context) {
	expertID := c.Param("expert_id")
	if expertID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Expert ID is required"})
		return
	}
	offTimes, err := h.offTimeService.GetExpertOffTimes(expertID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offTimes)
}

// Lấy chi tiết off-time theo id
func (h *OffTimeHandler) GetOffTimeByID(c *gin.Context) {
	id := c.Param("id")
	offTime, err := h.offTimeService.GetOffTimeByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offTime)
}

// Xóa off-time theo id
func (h *OffTimeHandler) DeleteOffTime(c *gin.Context) {
	id := c.Param("id")
	if err := h.offTimeService.DeleteOffTime(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
