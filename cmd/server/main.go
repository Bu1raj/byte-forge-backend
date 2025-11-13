package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bu1raj/byte-forge-backend/internal/api"
	"github.com/Bu1raj/byte-forge-backend/internal/store/redis"
)

var redisConfig = &redis.RedisStoreConfig{
	Address:  "localhost:6379",
	Password: "",
	DB:       0,
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := api.NewServer(redisConfig)

	server.Start(ctx)

	<-ctx.Done()
	log.Println("shutting down the server gracefully...")

	server.Shutdown(ctx)
}
