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
	"github.com/Bu1raj/byte-forge-backend/internal/store"
)

// TODO this should come from vault or env vars
var config = &store.KafkaStoreConfig{
	Broker:         "localhost:9092",
	ProducerTopics: []string{"submissions"},
	ConsumerTopics: []string{"results"},
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// Initialize Kafka producers and start the consumers
	store.InitKafkaUtilStore(config)
	background.StartResultConsumer(ctx, &wg)

	mux := http.NewServeMux()
	RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("listening :8080")

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
