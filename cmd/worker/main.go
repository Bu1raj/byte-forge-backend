package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	judgezero "github.com/Bu1raj/byte-forge-backend/internal/judge-zero"
	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"github.com/Bu1raj/byte-forge-backend/internal/store"
	kafkaStore "github.com/Bu1raj/byte-forge-backend/internal/store/kafka"
	kafka "github.com/segmentio/kafka-go"
)

// need to store these in vault
var kafkaConfig = &kafkaStore.KafkaStoreConfig{
	Broker:         "localhost:29092",
	ProducerTopics: []string{"results"},
	ConsumerTopics: []string{"submissions"},
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store := store.InitStore(kafkaConfig, nil)

	handleCodeSubmissions := func(msg *kafka.Message) error {
		var job models.KafkaCodeSubmissionsPayload
		err := json.Unmarshal(msg.Value, &job)
		if err != nil {
			return fmt.Errorf("failed to unmarshal the kafka message: %v", err)
		}

		log.Printf("Processing job %s", job.ID)
		res, err := judgezero.SubmitCode(job.Request)
		if err != nil {
			log.Printf("failed to submit code to judge0: %v", err)
			return err
		}

		result := models.KafkaCodeResultsPayload{
			ID:     job.ID,
			Result: res,
		}
		// Publish results to results topic
		data, _ := json.Marshal(result)
		results_producer, ok := store.Kafka.GetProducer("results")
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

	submissions_consumer, ok := store.Kafka.GetConsumer("submissions")
	if !ok {
		log.Fatalf("Submissions consumer not found in store")
	}

	log.Println("Starting submission consumer...")
	err := submissions_consumer.Consume(ctx, handleCodeSubmissions)
	if err != nil {
		log.Printf("submissions consumer exited with error: %v", err)
	}
}
