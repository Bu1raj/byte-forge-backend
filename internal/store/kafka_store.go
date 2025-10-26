package store

import (
	"log"
	"os"
	"sync"

	"github.com/Bu1raj/byte-forge-backend/internal/queue"
)

type KafkaStoreConfig struct {
	Broker         string
	ProducerTopics []string
	ConsumerTopics []string
}

type kafkaUtilStore struct {
	broker    string
	producers map[string]*queue.Producer
	consumers map[string]*queue.Consumer
	mu        sync.Mutex
}

var (
	globalKafkaUtilStore *kafkaUtilStore
)

// LoadKafkaConfigFromEnv loads Kafka configuration from environment variables
func LoadKafkaConfigFromEnv(producerDefaults, consumerDefaults []string) *KafkaStoreConfig {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:29092"
	}

	producerTopics := getEnvAsSlice("KAFKA_PRODUCER_TOPICS", producerDefaults)
	consumerTopics := getEnvAsSlice("KAFKA_CONSUMER_TOPICS", consumerDefaults)

	return &KafkaStoreConfig{
		Broker:         broker,
		ProducerTopics: producerTopics,
		ConsumerTopics: consumerTopics,
	}
}

// InitKafkaUtilStoreFromEnv initializes Kafka store from environment variables
func InitKafkaUtilStoreFromEnv(producerDefaults, consumerDefaults []string) {
	config := LoadKafkaConfigFromEnv(producerDefaults, consumerDefaults)
	InitKafkaUtilStore(config)
}

// Init initializes all Kafka producers/consumers once at startup
func InitKafkaUtilStore(config *KafkaStoreConfig) {
	log.Println("[store] initializing kafka store...")

	globalKafkaUtilStore = &kafkaUtilStore{
		broker:    config.Broker,
		producers: make(map[string]*queue.Producer),
		consumers: make(map[string]*queue.Consumer),
	}

	for _, topic := range config.ProducerTopics {
		globalKafkaUtilStore.producers[topic] = queue.NewProducer(config.Broker, topic)
	}
	for _, topic := range config.ConsumerTopics {
		globalKafkaUtilStore.consumers[topic] = queue.NewConsumer(config.Broker, topic, topic+"-group")
	}
	log.Println("[store] kafka store initialized")
}

// getEnvAsSlice gets an environment variable as a comma-separated slice
func getEnvAsSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	// Simple split by comma
	var result []string
	for _, v := range splitByComma(value) {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

// splitByComma splits a string by comma
func splitByComma(s string) []string {
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

// TODO Can add methods to register new producers/consumers dynamically if needed
// Can also add methods to unregister/close the producers/consumers

// GetProducer retrieves the producer for the given topic.
// It returns the producer and a boolean indicating if it exists.
func GetProducer(topic string) (*queue.Producer, bool) {
	globalKafkaUtilStore.mu.Lock()
	defer globalKafkaUtilStore.mu.Unlock()
	producer, exists := globalKafkaUtilStore.producers[topic]
	return producer, exists
}

// GetConsumer retrieves the consumer for the given topic.
// It returns the consumer and a boolean indicating if it exists.
func GetConsumer(topic string) (*queue.Consumer, bool) {
	globalKafkaUtilStore.mu.Lock()
	defer globalKafkaUtilStore.mu.Unlock()
	consumer, exists := globalKafkaUtilStore.consumers[topic]
	return consumer, exists
}

// CloseAll closes all registered producers and consumers.
func CloseAll() {
	globalKafkaUtilStore.mu.Lock()
	defer globalKafkaUtilStore.mu.Unlock()
	for _, producer := range globalKafkaUtilStore.producers {
		producer.Close()
	}
	for _, consumer := range globalKafkaUtilStore.consumers {
		consumer.Close()
	}
}
