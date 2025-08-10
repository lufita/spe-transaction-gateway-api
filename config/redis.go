package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	lookup "spe-trx-gateway/lookup"
	"time"
)

func contextRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: lookup.RedisHost + ":" + lookup.RedisPort,
	})
	return client
}

func WriteRedis(key string, body string, expiry time.Duration) (string, error) {
	ctx := context.Background()
	client := contextRedis()
	defer client.Close()
	return client.Set(ctx, key, body, expiry).Result()
}

func ReadRedis(key string) (string, error) {
	ctx := context.Background()
	client := contextRedis()
	resultMessage, err := client.Get(ctx, key).Result()
	defer client.Close()
	return resultMessage, err
}
