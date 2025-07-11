// services/expert-service/internal/utils/service_token.go
package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// ServiceTokenClaims represents claims for service-to-service authentication
type ServiceTokenClaims struct {
	Service string `json:"service"`
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateServiceToken tạo JWT token cho service-to-service calls
func GenerateServiceToken(serviceName string) (string, error) {
	// Tạo một service user ID
	serviceUserID := uuid.New().String()

	// Tạo claims cho service
	claims := ServiceTokenClaims{
		Service: serviceName,
		UserID:  serviceUserID,
		Role:    "service",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Token valid for 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    serviceName,
			Subject:   serviceUserID,
		},
	}

	// Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token với secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret" // fallback, nhưng nên set trong env
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
