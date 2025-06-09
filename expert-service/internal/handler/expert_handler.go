package handler

import (
    "expert-service/internal/model"
    "expert-service/internal/service"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// ExpertHandler handles HTTP requests for expert resources.
type ExpertHandler struct {
    expertService service.ExpertService
}

// NewExpertHandler creates a new ExpertHandler.
func NewExpertHandler(expertService service.ExpertService) *ExpertHandler {
    return &ExpertHandler{
        expertService: expertService,
    }
}

// CreateExpert handles expert creation.
func (h *ExpertHandler) CreateExpert(c *gin.Context) {
    var req model.CreateExpertRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    expert, err := h.expertService.CreateExpert(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expert: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, expert)
}

// GetExpert returns expert by ID.
func (h *ExpertHandler) GetExpert(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expert ID"})
        return
    }

    expert, err := h.expertService.GetExpertByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Expert not found"})
        return
    }

    c.JSON(http.StatusOK, expert)
}

// GetExperts returns a paginated list of experts.
func (h *ExpertHandler) GetExperts(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "10")
    offsetStr := c.DefaultQuery("offset", "0")

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
        return
    }

    offset, err := strconv.Atoi(offsetStr)
    if err != nil || offset < 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
        return
    }

    experts, err := h.expertService.GetAllExperts(limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get experts: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "experts": experts,
        "limit":   limit,
        "offset":  offset,
    })
}

// UpdateExpert updates expert info by ID.
func (h *ExpertHandler) UpdateExpert(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expert ID"})
        return
    }

    var req model.UpdateExpertRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    if err := h.expertService.UpdateExpert(id, &req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expert: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Expert updated successfully"})
}

// DeleteExpert deletes expert by ID.
func (h *ExpertHandler) DeleteExpert(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expert ID"})
        return
    }

    if err := h.expertService.DeleteExpert(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expert: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Expert deleted successfully"})
}

// GetExpertsByExpertise returns experts by expertise.
func (h *ExpertHandler) GetExpertsByExpertise(c *gin.Context) {
    expertise := c.Query("expertise")
    if expertise == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Expertise parameter is required"})
        return
    }

    experts, err := h.expertService.GetExpertsByExpertise(expertise)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get experts: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"experts": experts})
}