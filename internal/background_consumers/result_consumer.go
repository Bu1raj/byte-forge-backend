package backgroundconsumers

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"github.com/Bu1raj/byte-forge-backend/internal/store"
	"github.com/segmentio/kafka-go"
)

type ResultConsumer struct {
	store *store.Store
}

// NewResultConsumer creates a new ResultConsumer instance
func NewResultConsumer(store *store.Store) *ResultConsumer {
	return &ResultConsumer{
		store: store,
	}
}

// StartResultConsumer starts a background consumer to listen for code execution results
func (rc *ResultConsumer) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	resultConsumer, ok := rc.store.Kafka.GetConsumer("results")
	if !ok {
		log.Println("Results consumer not found in store")
		// TODO Register a new consumer in the store
		return
	}

	log.Printf("Starting results consumer on topic: results")

	ResultConsumerHandler := func(msg *kafka.Message) error {
		var payload models.KafkaCodeResultsPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("Failed to unmarshal result payload: %v", err)
			return err
		}

		// Store the result in Redis
		err := rc.store.Redis.Store(ctx, payload.ID, payload.Result)
		if err != nil {
			log.Printf("Failed to store result for %s: %v", payload.ID, err)
			return err
		}

		log.Printf("Stored result for %s", payload.ID)
		return nil
	}

	err := resultConsumer.Consume(ctx, ResultConsumerHandler)
	if err != nil {
		log.Printf("Result consumer stopped: %v", err)
	}
}
