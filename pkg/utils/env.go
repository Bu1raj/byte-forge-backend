package utils

import (
	"fmt"
	"os"
)

// GetEnv gets an environment variable or returns a default value
func GetEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt gets an environment variable as an integer
func GetEnvAsInt(key string, defaultValue int) int {
	if value := GetEnv(key, ""); value != "" {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvAsSlice gets an environment variable as a comma-separated slice
func GetEnvAsSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var result []string
	for _, v := range SplitByComma(value) {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}
