package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"services/user-service/middleware"
	"services/user-service/model"
	"services/user-service/service"
	"services/user-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterRoutes(r *gin.Engine, userService service.UserService, jwtService service.JWTService) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", Register(userService))
		userGroup.POST("/login", Login(userService, jwtService))
		userGroup.POST("/refresh", RefreshToken(jwtService))
		userGroup.GET("/profile", middleware.AuthMiddleware(jwtService), GetProfile(userService))
		userGroup.PUT("/profile", middleware.AuthMiddleware(jwtService), UpdateProfile(userService))
		userGroup.DELETE("/profile", middleware.AuthMiddleware(jwtService), DeleteProfile(userService))
		userGroup.GET("/bookings", middleware.AuthMiddleware(jwtService), GetBookingHistory())
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

func RefreshToken(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken := c.GetHeader("X-Refresh-Token")
		if refreshToken == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
			return
		}

		claims, err := jwtService.ValidateToken(refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in token"})
			return
		}

		dummyUser := &model.User{ID: userID}

		newAccessToken, newRefreshToken, err := jwtService.GenerateTokens(dummyUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		})
	}
}

func GetProfile(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
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
		userIDStr := c.GetString("userID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
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
		userIDStr := c.GetString("userID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		if err := userService.DeleteProfile(userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func GetBookingHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")
		status := c.Query("status")
		// ... các query khác nếu cần

		// Build URL booking-service
		bookingServiceURL := os.Getenv("BOOKING_SERVICE_URL") // ví dụ: http://booking-service:8082
		url := fmt.Sprintf("%s/GetUserBookings?user_id=%s&page=%s&limit=%s", bookingServiceURL, userID, page, limit)
		if status != "" {
			url += "&status=" + status
		}

		// Gửi request sang booking-service
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}
		// Truyền JWT từ header gốc
		req.Header.Set("Authorization", c.GetHeader("Authorization"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot connect to booking service"})
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", body)
	}
}
