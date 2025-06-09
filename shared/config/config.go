package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// Server
	Port string
	Host string
	Env  string

	// Database
	DatabaseURL      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	// Redis
	RedisURL      string
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// JWT
	JWTSecret          string
	JWTExpiry          int
	RefreshTokenExpiry int

	// Email
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string

	// External APIs
	TelegramBotToken string
	AblyAPIKey       string
}

func LoadConfig() *Config {
	return &Config{
		// Server
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", "localhost"),
		Env:  getEnv("ENV", "development"),

		// Database
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		DatabaseHost:     getEnv("DB_HOST", "localhost"),
		DatabasePort:     getEnv("DB_PORT", "5432"),
		DatabaseUser:     getEnv("DB_USER", "postgres"),
		DatabasePassword: getEnv("DB_PASSWORD", "password"),
		DatabaseName:     getEnv("DB_NAME", "consultation_booking"),

		// Redis
		RedisURL:      getEnv("REDIS_URL", ""),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// JWT
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-key"),
		JWTExpiry:          getEnvInt("JWT_EXPIRY", 3600),       // 1 hour
		RefreshTokenExpiry: getEnvInt("REFRESH_EXPIRY", 604800), // 7 days

		// Email
		SMTPHost:     getEnv("SMTP_HOST", "localhost"),
		SMTPPort:     getEnvInt("SMTP_PORT", 1025),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),

		// External APIs
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		AblyAPIKey:       getEnv("ABLY_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}