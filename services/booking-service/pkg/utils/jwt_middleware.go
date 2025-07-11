package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		// Log để debug
		log.Printf("Received token: %s", tokenString)
		log.Printf("JWT Secret: %s", secret)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Đảm bảo token sử dụng signing method chính xác
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			log.Printf("Token parsing error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			log.Printf("Token is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Invalid token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Log claims để debug
		log.Printf("Token claims: %+v", claims)

		// Kiểm tra xem có phải service call không
		if service, exists := claims["service"].(string); exists {
			log.Printf("Service call from: %s", service)
			c.Set("service_call", true)
			c.Set("service_name", service)
			c.Set("user_role", "service") // Add this line to set role for service calls

			// Đối với service call, set user_id từ claims
			if userID, exists := claims["user_id"].(string); exists {
				c.Set("user_id", userID)
			}

			c.Next()
			return
		}

		// Xử lý normal user token
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			userIDStr, ok = claims["sub"].(string)
			if !ok {
				log.Printf("user_id missing in token")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id missing in token"})
				return
			}
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			log.Printf("user_id invalid format: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id invalid format"})
			return
		}

		c.Set("user_id", userID)

		// Thêm đoạn này để set user_role
		role, ok := claims["role"].(string)
		if ok {
			c.Set("user_role", role)
		}

		c.Next()
	}
}
