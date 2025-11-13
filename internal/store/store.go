package store

import (
	"github.com/Bu1raj/byte-forge-backend/internal/store/kafka"
	"github.com/Bu1raj/byte-forge-backend/internal/store/redis"
)

type Store struct {
	Kafka *kafka.KafkaUtilStore
	Redis *redis.RedisStore
}

// InitStore initializes the Store with Kafka and Redis stores
func InitStore(redisConfig *redis.RedisStoreConfig) *Store {
	kafkaStore := kafka.NewKafkaUtilStore()

	var redisStore *redis.RedisStore
	if redisConfig != nil {
		redisStore = redis.NewRedisStore(redisConfig)
	}

	return &Store{
		Kafka: kafkaStore,
		Redis: redisStore,
	}
}

// CloseStore closes all resources in the Store
func (st *Store) CloseStore() {
	if st.Kafka != nil {
		st.Kafka.CloseAll()
	}
	if st.Redis != nil {
		st.Redis.Close()
	}
}
