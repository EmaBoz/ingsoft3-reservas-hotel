package config

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ() *amqp.Connection {
	rabbitURL := getenvOrDefault("RABBITMQ_URL", "amqp://rabbitmq:5672/")

	var (
		conn *amqp.Connection
		err  error
	)

	// Backoff exponencial para conectar a RabbitMQ
	for attempt := 1; attempt <= 10; attempt++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			log.Printf("Conectado a RabbitMQ en %s (intento %d)", rabbitURL, attempt)
			return conn
		}

		wait := time.Duration(attempt*2) * time.Second
		log.Printf("RabbitMQ no listo (intento %d): %v. Reintentando en %s...", attempt, err, wait)
		time.Sleep(wait)
	}

	log.Fatalf("No se pudo conectar a RabbitMQ tras reintentos: %v", err)
	return nil
}
