package service

import (
	"encoding/json"
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

func handleUserDeletedEvent(d amqp.Delivery, productService *ProductService) {
	type EventStruct struct {
		Message string `json:"message"`
		UserID  uint   `json:"userID"`
	}
	var event EventStruct
	err := json.Unmarshal(d.Body, &event)
	if err != nil {
		log.Printf("Error unmarshalling event: %v", err)
		return
	}
	userID := event.UserID

	log.Printf("%v", event)
	err = productService.DeleteProductsByUserID(userID)
	if err != nil {
		log.Printf("Error deleting products for user %d: %v", userID, err)
	} else {
		log.Printf("Deleted products for user %d", userID)
	}

	// Acknowledge the message so it's not redelivered
	if err := d.Ack(false); err != nil {
		log.Printf("Error acknowledging message: %v", err)
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

func (r *RabbitMQService) ConsumeEvents(queueName, consumerTag string, handleMsg func(amqp.Delivery)) {
	ch, err := r.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			log.Fatalf("Failed to close channel: %v", err)
		}
	}()

	msgs, err := ch.Consume(
		queueName,   // queue
		consumerTag, // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handleMsg(d) // Process each message
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (r *RabbitMQService) StartListeningForUserDeleteEvents(productService *ProductService) {
	queueName := "userDeleteQueue" // Ensure this matches the queue where user delete events are published
	consumerTag := "productServiceConsumer"

	handleMsgWrapper := func(d amqp.Delivery) {
		handleUserDeletedEvent(d, productService)
	}

	r.ConsumeEvents(queueName, consumerTag, handleMsgWrapper)
}
