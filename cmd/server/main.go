package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Bu1raj/byte-forge-backend/internal/background"
	"github.com/Bu1raj/byte-forge-backend/internal/config"
	"github.com/Bu1raj/byte-forge-backend/internal/store"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// Load configuration from environment variables
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize Kafka producers and start the consumers
	kafkaConfig := &store.KafkaStoreConfig{
		Broker:         cfg.Kafka.Broker,
		ProducerTopics: cfg.Kafka.ProducerTopics,
		ConsumerTopics: cfg.Kafka.ConsumerTopics,
	}
	store.InitKafkaUtilStore(kafkaConfig)
	background.StartResultConsumer(ctx, &wg)

	mux := http.NewServeMux()
	RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	// Start the HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("listening :%s", cfg.Server.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down the server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	log.Println("waiting for background tasks to finish...")
	wg.Wait()
	log.Println("all background tasks finished")

	log.Println("closing all kafka connections...")
	store.CloseAll()
	log.Println("all kafka connections closed")
	log.Println("server stopped")
}
