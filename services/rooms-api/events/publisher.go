package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"rooms-api/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewEventPublisher(conn *amqp.Connection) (*EventPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declarar exchange "rooms" para search-api
	err = ch.ExchangeDeclare(
		"rooms",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &EventPublisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *EventPublisher) PublishRoomCreated(room *domain.Room) error {
	event := map[string]interface{}{
		"event_type": "created",
		"room_id":    room.ID,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	return p.publish("room.created", event)
}

func (p *EventPublisher) PublishRoomUpdated(room *domain.Room) error {
	event := map[string]interface{}{
		"event_type": "updated",
		"room_id":    room.ID,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	return p.publish("room.updated", event)
}

func (p *EventPublisher) PublishRoomDeleted(roomID uint) error {
	event := map[string]interface{}{
		"event_type": "deleted",
		"room_id":    roomID,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	return p.publish("room.deleted", event)
}

func (p *EventPublisher) PublishRoomStatusChanged(roomID uint, oldStatus, newStatus string) error {
	// Cambio de estado también es una actualización, usar routing key room.updated
	event := map[string]interface{}{
		"event_type": "updated",
		"room_id":    roomID,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	return p.publish("room.updated", event)
}

func (p *EventPublisher) publish(routingKey string, event map[string]interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	err = p.channel.Publish(
		"rooms",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	log.Printf("✅ Evento publicado: %s", routingKey)
	return nil
}

func (p *EventPublisher) Close() error {
	if p.channel != nil {
		return p.channel.Close()
	}
	return nil
}
