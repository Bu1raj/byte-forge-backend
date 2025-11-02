package config

import (
	"github.com/Bu1raj/byte-forge-backend/pkg/utils"
)

// GetServerPort returns the server port from environment or default
func GetServerPort() string {
	return utils.GetEnv("SERVER_PORT", "8080")
}
