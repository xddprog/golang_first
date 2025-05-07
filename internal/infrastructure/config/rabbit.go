package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)


type RabbitConfig struct {
	Host string
	Port int
	User string
	Password string 
}


func (cfg *RabbitConfig) ConnectionString() string {
	return cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + strconv.Itoa(cfg.Port)
}


func LoadRabbitConfig() (*RabbitConfig) {
	godotenv.Load()

	port, err := strconv.Atoi(os.Getenv("RABBIT_PORT"))
	if err != nil {
		panic("failed to parse RABBIT_PORT")
	}

	return &RabbitConfig{
		Host: os.Getenv("RABBIT_HOST"),
		Port: port,
		User: os.Getenv("RABBIT_USER"),
		Password: os.Getenv("RABBIT_PASS"),
	}
}