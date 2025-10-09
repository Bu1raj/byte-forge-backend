package background

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"github.com/Bu1raj/byte-forge-backend/internal/store"
	"github.com/segmentio/kafka-go"
)

// StartResultConsumer starts a background consumer to listen for code execution results
func StartResultConsumer(ctx context.Context, wg *sync.WaitGroup) {
	result_consumer, ok := store.GetConsumer("results")
	if !ok {
		log.Println("Results consumer not found in store")
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Starting results consumer on topic: results")
		err := result_consumer.Consume(ctx, ResultConsumerHandler)
		if err != nil {
			log.Printf("Result consumer stopped: %v", err)
		}
	}()
}

// ResultConsumerHandler processes messages from the results topic and stores them.
func ResultConsumerHandler(msg *kafka.Message) error {
	var payload models.KafkaCodeResultsPayload
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		log.Printf("Failed to unmarshal result payload: %v", err)
		return err
	}
	err := store.StoreResult(payload.ID, payload.Result)
	if err != nil {
		log.Printf("Failed to store result for %s: %v", payload.ID, err)
		return err
	}
	log.Printf("Stored result for %s", payload.ID)
	return nil
}
