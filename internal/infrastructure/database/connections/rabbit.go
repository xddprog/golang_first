package connections

import (
	"golang/internal/infrastructure/config"

	"github.com/rabbitmq/amqp091-go"
)


func NewRabbitConnection() (*amqp091.Connection, error) {
	rabbitCfg := config.LoadRabbitConfig()
	conn, err := amqp091.Dial(rabbitCfg.ConnectionString())
	if err != nil {
		return nil, err
	}
	return conn, nil
}