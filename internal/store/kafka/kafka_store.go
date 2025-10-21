package kafka

import (
	"log"

	"github.com/Bu1raj/byte-forge-backend/internal/queue"
)

type KafkaStoreConfig struct {
	Broker         string
	ProducerTopics []string
	ConsumerTopics []string
}

type KafkaUtilStore struct {
	broker    string
	producers map[string]*queue.Producer
	consumers map[string]*queue.Consumer
}

// NewKafkaUtilStore initializes all Kafka producers/consumers once at startup
func NewKafkaUtilStore(config *KafkaStoreConfig) *KafkaUtilStore {
	log.Println("[store] initializing kafka store...")

	kafkaStore := &KafkaUtilStore{
		broker:    config.Broker,
		producers: make(map[string]*queue.Producer),
		consumers: make(map[string]*queue.Consumer),
	}

	for _, topic := range config.ProducerTopics {
		kafkaStore.producers[topic] = queue.NewProducer(config.Broker, topic)
	}
	for _, topic := range config.ConsumerTopics {
		kafkaStore.consumers[topic] = queue.NewConsumer(config.Broker, topic, topic+"-group")
	}
	log.Println("[store] kafka store initialized")
	return kafkaStore
}

// TODO Can add methods to register new producers/consumers dynamically if needed
// in this case, need to add mutex to protect the maps
// Can also add methods to unregister/close the producers/consumers

// GetProducer retrieves the producer for the given topic.
// It returns the producer and a boolean indicating if it exists.
func (k *KafkaUtilStore) GetProducer(topic string) (*queue.Producer, bool) {
	producer, exists := k.producers[topic]
	return producer, exists
}

// GetConsumer retrieves the consumer for the given topic.
// It returns the consumer and a boolean indicating if it exists.
func (k *KafkaUtilStore) GetConsumer(topic string) (*queue.Consumer, bool) {
	consumer, exists := k.consumers[topic]
	return consumer, exists
}

// CloseAll closes all registered producers and consumers.
func (k *KafkaUtilStore) CloseAll() {
	for _, producer := range k.producers {
		producer.Close()
	}
	for _, consumer := range k.consumers {
		consumer.Close()
	}
}
