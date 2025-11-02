package config

import "fmt"

// GetServerPort returns the server port from environment or default
func GetServerPort() string {
	return GetEnv("SERVER_PORT", "8080")
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
