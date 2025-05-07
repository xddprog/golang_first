package connections

import (
	"golang/internal/infrastructure/config"

	"github.com/redis/go-redis/v9"
)


func NewRedisConnection() *redis.Client {
	redisCfg := config.LoadRedisConfig()
	return redis.NewClient(&redis.Options{
		Addr: redisCfg.GetAddress(),
		DB: redisCfg.DB,
		PoolSize: redisCfg.PoolSize,
	})
}