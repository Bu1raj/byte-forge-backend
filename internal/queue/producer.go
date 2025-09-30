package queue

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer, by initializing a kafka.Writer
// with the provided broker address and topic.
func NewProducer(broker string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// SendMessage sends a message to the Kafka topic.
func (p *Producer) SendMessage(msg []byte) error {
	// a defensive timeout context, in case Kafka is not reachable
	// instead of hanging forever, we fail after 10 seconds
	writeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := p.writer.WriteMessages(writeCtx,
		kafka.Message{
			Key:   []byte(time.Now().Format(time.RFC3339)),
			Value: msg,
		},
	)

	if err != nil {
		log.Printf("failed to write message: %v", err)
		return err
	}

	return nil
}
