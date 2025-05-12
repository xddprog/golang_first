package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)


type SmtpConfig struct {
	Host     string 
	Port     int    
	User     string 
	Password string 
}


func LoadSmtpConfig() *SmtpConfig {
	godotenv.Load()
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(err)
	}
	return &SmtpConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		User:     os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}