package config

import (
	"os"
)

type Config struct {
	Port        string
	BookingURL  string
	ExpertURL   string
	NotifyURL   string
}

func LoadConfig() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		BookingURL: getEnv("BOOKING_URL", "http://localhost:9001"),
		ExpertURL:  getEnv("EXPERT_URL", "http://localhost:9002"),
		NotifyURL:  getEnv("NOTIFY_URL", "http://localhost:9003"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
