package config

import (
	"os"
)

type Config struct {
	UserURL      string
	BookingURL   string
	ExpertURL    string
	NotifyURL    string
	JWTSecret    string
	RateLimit    int
	RateDuration int
}

func NewConfig() *Config {
	return &Config{
		UserURL:      getEnv("USER_SERVICE_URL", "http://localhost:8080"),
		BookingURL:   getEnv("BOOKING_SERVICE_URL", "http://localhost:8081"),
		ExpertURL:    getEnv("EXPERT_SERVICE_URL", "http://localhost:8082"),
		NotifyURL:    getEnv("NOTIFY_SERVICE_URL", "http://localhost:8083"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		RateLimit:    100,
		RateDuration: 60,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
