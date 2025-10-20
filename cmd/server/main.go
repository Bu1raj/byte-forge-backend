package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bu1raj/byte-forge-backend/internal/api"
	"github.com/Bu1raj/byte-forge-backend/internal/store/kafka"
	"github.com/Bu1raj/byte-forge-backend/internal/store/redis"
)

// TODO this should come from vault or env vars
var kafkaConfig = &kafka.KafkaStoreConfig{
	Broker:         "localhost:9092",
	ProducerTopics: []string{"submissions"},
	ConsumerTopics: []string{"results"},
}

var redisConfig = &redis.RedisStoreConfig{
	Address:  "localhost:6379",
	Password: "",
	DB:       0,
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := api.NewServer(kafkaConfig, redisConfig)

	server.Start(ctx)

	<-ctx.Done()
	log.Println("shutting down the server gracefully...")

	server.Shutdown(ctx)
}
