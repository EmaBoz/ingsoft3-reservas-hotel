package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reservations-api/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "reservations"
	ExchangeType = "topic"
)

type EventPublisher interface {
	PublishReservationCreated(ctx context.Context, event domain.ReservationEvent) error
	PublishReservationCanceled(ctx context.Context, event domain.ReservationEvent) error
	Close() error
}

type rabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewEventPublisher(conn *amqp.Connection) (EventPublisher, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declarar el exchange
	err = channel.ExchangeDeclare(
		ExchangeName, // name
		ExchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		channel.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	log.Printf("Exchange '%s' declarado exitosamente", ExchangeName)

	return &rabbitMQPublisher{
		conn:    conn,
		channel: channel,
	}, nil
}

func (p *rabbitMQPublisher) PublishReservationCreated(ctx context.Context, event domain.ReservationEvent) error {
	return p.publish(ctx, "reservation.created", event)
}

func (p *rabbitMQPublisher) PublishReservationCanceled(ctx context.Context, event domain.ReservationEvent) error {
	return p.publish(ctx, "reservation.canceled", event)
}

func (p *rabbitMQPublisher) publish(ctx context.Context, routingKey string, event domain.ReservationEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		ExchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Evento publicado: %s - Reservation ID: %s", routingKey, event.ReservationID)
	return nil
}

func (p *rabbitMQPublisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return err
		}
	}
	return nil
}
