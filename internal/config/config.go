package config

import (
	"github.com/Bu1raj/byte-forge-backend/pkg/utils"
)

type KafkaStoreConfig struct {
	Broker         string
	ProducerTopics []string
	ConsumerTopics []string
}

type RedisStoreConfig struct {
	Address  string
	Password string
	DB       int
}

// GetServerPort returns the server port from environment or default
func GetServerPort() string {
	return utils.GetEnv("SERVER_PORT", "8080")
}

// GetKafkaConfig returns the Kafka store configuration from environment variables
func GetKafkaConfig() *KafkaStoreConfig {
	return &KafkaStoreConfig{
		Broker:         utils.GetEnv("KAFKA_BROKER", "localhost:29092"),
		ProducerTopics: utils.GetEnvAsSlice("KAFKA_PRODUCER_TOPICS", []string{"submissions"}),
		ConsumerTopics: utils.GetEnvAsSlice("KAFKA_CONSUMER_TOPICS", []string{"results"}),
	}
}

// GetRedisConfig returns the Redis store configuration from environment variables
func GetRedisConfig() *RedisStoreConfig {
	return &RedisStoreConfig{
		Address:  utils.GetEnv("REDIS_ADDRESS", "localhost:6379"),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       utils.GetEnvAsInt("REDIS_DB", 0),
	}
}
