package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	Kafka  KafkaConfig
	Redis  RedisConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string
}

// KafkaConfig holds Kafka-specific configuration
type KafkaConfig struct {
	Broker         string
	ProducerTopics []string
	ConsumerTopics []string
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Kafka: KafkaConfig{
			Broker:         getEnv("KAFKA_BROKER", "localhost:29092"),
			ProducerTopics: getEnvAsSlice("KAFKA_PRODUCER_TOPICS", []string{"submissions"}),
			ConsumerTopics: getEnvAsSlice("KAFKA_CONSUMER_TOPICS", []string{"results"}),
		},
		Redis: RedisConfig{
			Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

	return config, nil
}

// LoadServerConfig loads configuration for the server
func LoadServerConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Kafka: KafkaConfig{
			Broker:         getEnv("KAFKA_BROKER", "localhost:29092"),
			ProducerTopics: getEnvAsSlice("KAFKA_PRODUCER_TOPICS", []string{"submissions"}),
			ConsumerTopics: getEnvAsSlice("KAFKA_CONSUMER_TOPICS", []string{"results"}),
		},
		Redis: RedisConfig{
			Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

	return config, nil
}

// LoadWorkerConfig loads configuration for the worker
func LoadWorkerConfig() (*Config, error) {
	config := &Config{
		Kafka: KafkaConfig{
			Broker:         getEnv("KAFKA_BROKER", "localhost:29092"),
			ProducerTopics: getEnvAsSlice("KAFKA_PRODUCER_TOPICS", []string{"results"}),
			ConsumerTopics: getEnvAsSlice("KAFKA_CONSUMER_TOPICS", []string{"submissions"}),
		},
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsSlice gets an environment variable as a comma-separated slice
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}
