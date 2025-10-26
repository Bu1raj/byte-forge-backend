package config

import (
	"fmt"
	"os"
	"strings"
)

// GetServerPort returns the server port from environment or default
func GetServerPort() string {
	return getEnv("SERVER_PORT", "8080")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsSlice gets an environment variable as a comma-separated slice
func GetEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// GetEnvAsInt gets an environment variable as an integer
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}
