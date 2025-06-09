package service

import (
	"services/user-service/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateTokens(user *model.User) (string, string, error)
	ValidateToken(token string) (*jwt.RegisteredClaims, error)
}

type jwtService struct {
	secret             string
	expiry             int
	refreshTokenExpiry int
}

func NewJWTService(secret string, expiry int, refreshTokenExpiry int) JWTService {
	return &jwtService{secret, expiry, refreshTokenExpiry}
}

func (j *jwtService) GenerateTokens(user *model.User) (string, string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expiry) * time.Second)),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}
	refreshClaims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.refreshTokenExpiry) * time.Second)),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}
	return tokenString, refreshTokenString, nil
}

func (j *jwtService) ValidateToken(token string) (*jwt.RegisteredClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsed.Claims.(*jwt.RegisteredClaims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, err
}
