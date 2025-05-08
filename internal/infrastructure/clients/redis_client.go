package clients

import "github.com/redis/go-redis/v9"


type RedisClient struct {
	client *redis.Client
}


func NewRedisClient(host string, cli *redis.Client) *RedisClient {
	return &RedisClient{client: cli,}
}
