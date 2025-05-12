package clients

import (
	"golang/internal/infrastructure/config"

	"github.com/rabbitmq/amqp091-go"
)


type RabbitClient struct {
	Connection *amqp091.Connection
}


func NewRabbitClient() *RabbitClient {
	rabbitCfg := config.LoadRabbitConfig()
	conn, err := amqp091.Dial(rabbitCfg.ConnectionString())
	if err != nil {
		panic(err)
	}
	return &RabbitClient{Connection: conn}
}
