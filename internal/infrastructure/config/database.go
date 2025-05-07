package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)


type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}



func (cfg *DatabaseConfig) ConnectionString() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=disable",
        cfg.User,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.DBName,
    )
}


func LoadDatabaseConfig() (*DatabaseConfig, error) {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env file")
	}
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
	}, nil
}

