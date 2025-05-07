package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)


type RedisConfig struct {
	Host string
	Port int
	DB int
	PoolSize int
}


func (cfg *RedisConfig) GetAddress() string {
	return cfg.Host + ":" + strconv.Itoa(cfg.Port)
}


func LoadRedisConfig() (*RedisConfig) {
	godotenv.Load()

	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		panic("failed to parse REDIS_PORT")
	}

	return &RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
		Port: port,
		DB:   0,
		PoolSize: 50,
	}
}