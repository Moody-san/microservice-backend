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
	// Open a channel
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close() // Ensure the channel is closed after the function execution

	err = ch.ExchangeDeclare(
		exchange, // Use the dynamic exchange name
		"fanout", // Usually, fanout is used for broadcasting. If you need direct or topic, change accordingly.
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Publish the message
	err = ch.Publish(
		exchange,   // Use the dynamic exchange name
		routingKey, // Use the provided routing key. For fanout exchange, this can be ignored.
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			Headers:     amqp.Table{"eventType": eventType},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

func (r *RabbitMQService) SetupQueueAndBind(exchangeName, queueName string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Declare the fanout exchange
	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Declare the queue
	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		q.Name,
		"", // Ignored by fanout exchanges
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind the queue: %w", err)
	}

	return nil
}
