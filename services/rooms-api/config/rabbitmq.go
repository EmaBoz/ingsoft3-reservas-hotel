package config

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ(url string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	// Reintentar conexión hasta 10 veces
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			return conn, nil
		}

		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("no se pudo conectar a RabbitMQ después de 10 intentos: %w", err)
}
