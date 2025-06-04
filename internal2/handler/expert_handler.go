package handler

import (
    "expert-service/internal2/model"
    "expert-service/internal2/service"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

type ExpertHandler struct {
    expertService service.ExpertService
}

func NewExpertHandler(expertService service.ExpertService) *ExpertHandler {
    return &ExpertHandler{
        expertService: expertService,
    }
}

// CreateExpert tạo chuyên gia mới
func (h *ExpertHandler) CreateExpert(c *gin.Context) {
    var req model.CreateExpertRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Dữ liệu không hợp lệ",
            "details": err.Error(),
        })
        return
    }

    expert, err := h.expertService.CreateExpert(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Tạo chuyên gia thành công",
        "data":    expert,
    })
}

// GetExpert lấy thông tin chuyên gia theo ID
func (h *ExpertHandler) GetExpert(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "ID không hợp lệ",
        })
        return
    }

    expert, err := h.expertService.GetExpert(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": expert,
    })
}

// GetAllExperts lấy tất cả chuyên gia
func (h *ExpertHandler) GetAllExperts(c *gin.Context) {
    experts, err := h.expertService.GetAllExperts()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data":  experts,
        "count": len(experts),
    })
}

// UpdateExpert cập nhật thông tin chuyên gia
func (h *ExpertHandler) UpdateExpert(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "ID không hợp lệ",
        })
        return
    }

    var req model.CreateExpertRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Dữ liệu không hợp lệ",
            "details": err.Error(),
        })
        return
    }

    expert, err := h.expertService.UpdateExpert(id, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Cập nhật chuyên gia thành công",
        "data":    expert,
    })
}

// DeleteExpert xóa chuyên gia
func (h *ExpertHandler) DeleteExpert(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "ID không hợp lệ",
        })
        return
    }

    err = h.expertService.DeleteExpert(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Xóa chuyên gia thành công",
    })
}