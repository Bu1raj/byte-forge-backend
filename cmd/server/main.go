package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bu1raj/byte-forge-backend/internal/api"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := api.NewServer()

	server.Start(ctx)

	<-ctx.Done()
	log.Println("shutting down the server gracefully...")

	server.Shutdown(ctx)
}
