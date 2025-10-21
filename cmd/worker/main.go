package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bu1raj/byte-forge-backend/internal/config"
	"github.com/Bu1raj/byte-forge-backend/internal/executor"
	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"github.com/Bu1raj/byte-forge-backend/internal/store"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load configuration from environment variables
	cfg, err := config.LoadWorkerConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	kafkaConfig := &store.KafkaStoreConfig{
		Broker:         cfg.Kafka.Broker,
		ProducerTopics: cfg.Kafka.ProducerTopics,
		ConsumerTopics: cfg.Kafka.ConsumerTopics,
	}
	store.InitKafkaUtilStore(kafkaConfig)

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
