package clients

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)


type RedisClient struct {
	client *redis.Client
}


func NewRedisClient(host string, cli *redis.Client) *RedisClient {
	return &RedisClient{client: cli,}
}


func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	err := r.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}