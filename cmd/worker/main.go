package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bu1raj/byte-forge-backend/internal/executor"
	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"github.com/Bu1raj/byte-forge-backend/internal/store"
	"github.com/segmentio/kafka-go"
)

// need to store these in vault
var config = &store.KafkaStoreConfig{
	Broker:         "localhost:9092",
	ProducerTopics: []string{"results"},
	ConsumerTopics: []string{"submissions"},
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store.InitKafkaUtilStore(config)

	handleCodeSubmissions := func(msg *kafka.Message) error {
		var job models.KafkaCodeSubmissionsPayload
		err := json.Unmarshal(msg.Value, &job)
		if err != nil {
			return fmt.Errorf("failed to unmarshal the kafka message: %v", err)
		}

		log.Printf("Processing job %s", job.ID)

		res := executor.RunSubmission(job.ID, job.SubmitRequest)

		result := models.KafkaCodeResultsPayload{
			ID:     job.ID,
			Result: res,
		}
		// Publish results to results topic
		data, _ := json.Marshal(result)
		results_producer, ok := store.GetProducer("results")
		if !ok {
			log.Printf("Results producer not found in store")
			return fmt.Errorf("results producer not found in store")
		}
		err = results_producer.SendMessage(data)
		if err != nil {
			log.Printf("failed to publish results to kafka: %v", err)
		}

		fmt.Printf("Result: %+v\n", res)
		return nil
	}
	submissions_consumer, ok := store.GetConsumer("submissions")
	if !ok {
		log.Fatalf("Submissions consumer not found in store")
	}
	log.Println("Starting submission consumer...")
	submissions_consumer.Consume(ctx, handleCodeSubmissions)
}
