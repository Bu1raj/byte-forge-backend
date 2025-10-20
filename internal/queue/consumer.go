package queue

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer, by initializing a kafka.Reader
// with the provided broker address, topic, and group ID.
func NewConsumer(broker, topic, groupID string) *Consumer {
	log.Printf("Broker : %s\n", broker)
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

// Consume starts consuming messages from the Kafka topic.
// It takes a parent context to allow for graceful shutdown and
// a handler function to process each message.
func (c *Consumer) Consume(parentCtx context.Context, handler func(msg *kafka.Message) error) error {
	defer c.reader.Close()
	for {
		// Check if the parent context is done, faster shutdown
		select {
		case <-parentCtx.Done():
			log.Println("Parent context done, consumer shutting down")
			return parentCtx.Err()
		default:
		}

		readCtx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
		msg, err := c.reader.ReadMessage(readCtx)
		cancel() // explicitly call the cancel to cleanup the readCtx
		if err != nil {
			// this should occur when parent calls a cancel
			// parentCtx.Err() == context.Cancelled
			if errors.Is(err, context.Canceled) || errors.Is(parentCtx.Err(), context.Canceled) {
				log.Println("Received cancellation signal, consumer shutting down")
				return parentCtx.Err()
			}

			// read timeout (no messages in the last 30s)
			if readCtx.Err() == context.DeadlineExceeded {
				continue
			}

			log.Printf("unexpected consumer error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		// process the message
		if err := handler(&msg); err != nil {
			log.Printf("handler error: %v", err)
			// TODO do we want to requeue here ?
		}
	}
}

// Close closes the Kafka reader.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
