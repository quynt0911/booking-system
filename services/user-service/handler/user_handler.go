package handler

import (
	"net/http"
	"services/user-service/middleware"
	"services/user-service/model"
	"services/user-service/service"
	"services/user-service/utils"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userService service.UserService, jwtService service.JWTService) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", Register(userService))
		userGroup.POST("/login", Login(userService, jwtService))
		userGroup.GET("/profile", middleware.AuthMiddleware(jwtService), GetProfile(userService))
		userGroup.PUT("/profile", middleware.AuthMiddleware(jwtService), UpdateProfile(userService))
		userGroup.DELETE("/profile", middleware.AuthMiddleware(jwtService), DeleteProfile(userService))
	}
}

func Register(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.ParseValidationError(err)})
			return
		}
		if err := utils.ValidateRegisterInput(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := userService.Register(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func Login(userService service.UserService, jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.ParseValidationError(err)})
			return
		}
		user, err := userService.Login(req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		token, refreshToken, err := jwtService.GenerateTokens(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access_token": token, "refresh_token": refreshToken, "user": user})
	}
}

func GetProfile(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		user, err := userService.GetProfile(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func UpdateProfile(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		var req model.UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.ParseValidationError(err)})
			return
		}
		if err := utils.ValidateUpdateProfileInput(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := userService.UpdateProfile(userID, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func DeleteProfile(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if err := userService.DeleteProfile(userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
