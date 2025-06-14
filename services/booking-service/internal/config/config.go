package config

import (
	"os"
	"strconv"
)

type Config struct {
	App struct {
		Port        string
		Environment string
	}
	Database struct {
		URL string
	}
	Redis struct {
		URL string
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	// App config
	cfg.App.Port = getEnv("APP_PORT", "8082")
	cfg.App.Environment = getEnv("APP_ENV", "development")

	// Database config
	cfg.Database.URL = getEnv("BOOKING_SERVICE_DSN", getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/booking?sslmode=disable"))

	// Redis config
	cfg.Redis.URL = getEnv("REDIS_URL", "localhost:6379")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt reads an environment variable as an integer
// Returns the default value if the environment variable is not set or cannot be parsed as an integer
//
//nolint:unused
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
