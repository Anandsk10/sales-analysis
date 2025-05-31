package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	LogLevel    string
}

func New() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:root@localhost/salesdb?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
