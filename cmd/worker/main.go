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
	"github.com/Bu1raj/byte-forge-backend/internal/queue"
	"github.com/segmentio/kafka-go"
)

// need to store these in vault
const BROKER = "localhost:9092"
const TOPIC = "submissions"
const GROUP_ID = "submission-workers"

func main() {
	consumer := queue.NewConsumer(BROKER, TOPIC, GROUP_ID)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	handleCodeSubmissions := func(msg *kafka.Message) error {
		var job models.KafkaCodeSubmissionsPayload
		err := json.Unmarshal(msg.Value, &job)
		if err != nil {
			return fmt.Errorf("failed to unmarshal the kafka message: %v", err)
		}

		log.Printf("Processing job %s", job.ID)

		res := executor.RunSubmission(job.ID, job.SubmitRequest)
		// TODO we need publish the result to another topic
		// from there have the server consume and update the store
		// for now just printing the result
		fmt.Printf("Result: %+v\n", res)

		return nil
	}

	consumer.Consume(ctx, handleCodeSubmissions)
}
