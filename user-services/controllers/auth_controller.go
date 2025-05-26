package controllers

import (
	"user-services/config"
	"user-services/models"
	"user-services/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// Băm mật khẩu
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Email:    input.Email,
		Password: string(hashed),
		Name:     input.Name,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo người dùng"})
		return
	}

	go utils.SendWelcomeEmail(user.Email, user.Name)

	c.JSON(http.StatusCreated, gin.H{"message": "Đăng ký thành công"})
}
