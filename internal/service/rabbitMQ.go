package service

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn *amqp.Connection
}

func NewRabbitMQService(connStr string) *RabbitMQService {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	return &RabbitMQService{
		conn: conn,
	}
}

func (r *RabbitMQService) PublishEvent(exchange, routingKey, eventType string, payload []byte) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %w", err)
	}
	defer ch.Close()

	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			Headers:     amqp.Table{"eventType": eventType},
		})
	if err != nil {
		return fmt.Errorf("Failed to publish a message: %w", err)
	}

	return nil
}
