package store

import (
	"log"
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
