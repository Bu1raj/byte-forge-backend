package kafka

import (
	"log"

	"github.com/Bu1raj/byte-forge-backend/internal/config"
	"github.com/Bu1raj/byte-forge-backend/internal/queue"
)

type KafkaUtilStore struct {
	broker    string
	producers map[string]*queue.Producer
	consumers map[string]*queue.Consumer
}

// NewKafkaUtilStore initializes all Kafka producers/consumers once at startup
func NewKafkaUtilStore() *KafkaUtilStore {
	log.Println("[store] initializing kafka store...")

	kafkaConfig := config.GetKafkaConfig()

	kafkaStore := &KafkaUtilStore{
		broker:    kafkaConfig.Broker,
		producers: make(map[string]*queue.Producer),
		consumers: make(map[string]*queue.Consumer),
	}

	for _, topic := range kafkaConfig.ProducerTopics {
		log.Println("[store] initializing kafka producer for topic:", topic)
		kafkaStore.producers[topic] = queue.NewProducer(kafkaConfig.Broker, topic)
	}
	for _, topic := range kafkaConfig.ConsumerTopics {
		log.Println("[store] initializing kafka consumer for topic:", topic)
		kafkaStore.consumers[topic] = queue.NewConsumer(kafkaConfig.Broker, topic, topic+"-group")
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
