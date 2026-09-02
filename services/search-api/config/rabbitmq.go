package config

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConfig contiene la configuración de RabbitMQ
type RabbitMQConfig struct {
	URL          string
	ExchangeName string
	ExchangeType string
	QueueName    string
	RoutingKeys  []string
}

// NewRabbitMQConfig crea una nueva configuración de RabbitMQ
func NewRabbitMQConfig() *RabbitMQConfig {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	return &RabbitMQConfig{
		URL:          url,
		ExchangeName: "rooms",
		ExchangeType: "topic",
		QueueName:    "search-api-rooms-queue",
		RoutingKeys:  []string{"room.created", "room.updated", "room.deleted"},
	}
}

// ConnectWithRetry intenta conectarse a RabbitMQ con reintentos
func (c *RabbitMQConfig) ConnectWithRetry(maxRetries int, retryDelay time.Duration) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(c.URL)
		if err == nil {
			log.Printf("Successfully connected to RabbitMQ")
			return conn, nil
		}

		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
}

// SetupExchangeAndQueue configura el exchange y la queue
func (c *RabbitMQConfig) SetupExchangeAndQueue(ch *amqp.Channel) error {
	// Declarar exchange
	err := ch.ExchangeDeclare(
		c.ExchangeName, // name
		c.ExchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declarar queue
	queue, err := ch.QueueDeclare(
		c.QueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue al exchange con las routing keys
	for _, routingKey := range c.RoutingKeys {
		err = ch.QueueBind(
			queue.Name,     // queue name
			routingKey,     // routing key
			c.ExchangeName, // exchange
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue with routing key %s: %w", routingKey, err)
		}
		log.Printf("Queue %s bound to exchange %s with routing key %s", queue.Name, c.ExchangeName, routingKey)
	}

	return nil
}
