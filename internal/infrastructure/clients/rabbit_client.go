package clients

import "github.com/rabbitmq/amqp091-go"


type RabbitClient struct {
	Connection *amqp091.Connection
}


func NewRabbitClient(host string, connection *amqp091.Connection) *RabbitClient {
	return &RabbitClient{Connection: connection}
}
