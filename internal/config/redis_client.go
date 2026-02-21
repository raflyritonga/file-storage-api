package config

import (
	"context"

	redisv9 "github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg Config) (*redisv9.Client, error) {
	client := redisv9.NewClient(&redisv9.Options{
		Addr: cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB: cfg.RedisDB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}	

	return client, nil
}