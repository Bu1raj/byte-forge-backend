package api

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	bgConsumer "github.com/Bu1raj/byte-forge-backend/internal/background_consumers"
	"github.com/Bu1raj/byte-forge-backend/internal/store"
	"github.com/Bu1raj/byte-forge-backend/internal/store/redis"
)

type Server struct {
	HTTP    *http.Server
	Store   *store.Store
	Workers []bgConsumer.BgConsumer
	wg      *sync.WaitGroup
}

// NewServer creates a new Server instance
func NewServer(redisConfig *redis.RedisStoreConfig) *Server {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	var wg sync.WaitGroup
	store := store.InitStore(redisConfig)

	var backgroundConsumers []bgConsumer.BgConsumer
	resultConsumer := bgConsumer.NewResultConsumer(store)
	backgroundConsumers = append(backgroundConsumers, resultConsumer)

	return &Server{
		HTTP:    srv,
		Store:   store,
		Workers: backgroundConsumers,
		wg:      &wg,
	}
}

// Start starts the HTTP server and background workers
func (srv *Server) Start(ctx context.Context) {
	// register routes
	srv.RegisterRoutes(srv.HTTP.Handler.(*http.ServeMux))

	// start background processes
	for _, worker := range srv.Workers {
		srv.wg.Add(1)
		go worker.Run(ctx, srv.wg)
	}

	// start HTTP server
	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()

		log.Println("listening :8080")
		if err := srv.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

// Shutdown gracefully shuts down the server and background processes
func (srv *Server) Shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.HTTP.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	log.Println("waiting for background tasks to finish...")
	srv.wg.Wait()
	log.Println("all background tasks finished")

	log.Println("closing all kafka connections...")
	srv.Store.CloseStore()
	log.Println("all kafka connections closed")
	log.Println("server stopped")
}
