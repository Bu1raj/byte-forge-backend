package config

import "os"

// GetEnv gets an environment variable or returns a default value
func GetEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsSlice gets an environment variable as a comma-separated slice
func GetEnvAsSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	// Simple split by comma
	var result []string
	for _, v := range SplitByComma(value) {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

// SplitByComma splits a string by comma
func SplitByComma(s string) []string {
	var result []string
	current := ""
	for _, char := range s {
		if char == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
