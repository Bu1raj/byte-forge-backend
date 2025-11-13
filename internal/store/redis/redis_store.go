package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Bu1raj/byte-forge-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var ErrorKeyNotFound = errors.New("key not found in redis")

type RedisStore struct {
	client            *redis.Client
	defaultExpiration time.Duration
}

// NewRedisStore initializes a new RedisStore with configuration from environment variables.
func NewRedisStore() *RedisStore {
	redisConfig := config.GetRedisConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return nil
	}

	log.Println("Connected to Redis successfully")
	return &RedisStore{
		client:            client,
		defaultExpiration: 5 * time.Minute,
	}
}

// Store saves a value in Redis with the specified key and optional expiration time.
// If no expiration is provided, it uses the default expiration time which is 5 minutes.
// if expiration is 0, the key has no expiration time.
func (r *RedisStore) Store(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	expirationTime := r.defaultExpiration
	if len(expiration) > 0 {
		expirationTime = expiration[0]
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	if err := r.client.Set(ctx, key, data, expirationTime).Err(); err != nil {
		return fmt.Errorf("failed to store result in Redis: %v", err)
	}
	return nil
}

// Get retrieves a value from Redis by its key.
// NOTE: data should be a pointer to the expected type.
func (r *RedisStore) Get(ctx context.Context, key string, data interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrorKeyNotFound
	}
	if err != nil {
		log.Printf("Failed to get result from Redis: %v", err)
		return err
	}

	if err := json.Unmarshal([]byte(val), data); err != nil {
		log.Printf("Failed to unmarshal result from Redis: %v", err)
		return err
	}
	return nil
}

// Delete removes a key from Redis.
func (r *RedisStore) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key from Redis: %v", err)
	}
	return nil
}

// Close closes the Redis client connection.
func (r *RedisStore) Close() error {
	return r.client.Close()
}
