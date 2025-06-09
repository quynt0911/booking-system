package config

import (
    "os"
    "strconv"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
}

type ServerConfig struct {
    Port string
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type RedisConfig struct {
    Host     string
    Port     string
    Password string
    DB       int
}

func Load() *Config {
    dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
    redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

    return &Config{
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     dbPort,
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "password"),
            DBName:   getEnv("DB_NAME", "expert_service"),
            SSLMode:  getEnv("DB_SSL_MODE", "disable"),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnv("REDIS_PORT", "6379"),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       redisDB,
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}