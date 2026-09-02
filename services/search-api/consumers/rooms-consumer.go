package consumers

import (
	"encoding/json"
	"log"
	"time"

	"search-api/config"
	"search-api/domain"
	"search-api/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RoomsConsumer consume eventos de RabbitMQ para sincronizar Solr
type RoomsConsumer struct {
	rabbitMQConfig *config.RabbitMQConfig
	roomsAPIClient *config.RoomsAPIClient
	searchService  *services.SearchService
	conn           *amqp.Connection
	channel        *amqp.Channel
}

// NewRoomsConsumer crea un nuevo consumidor de eventos de rooms
func NewRoomsConsumer(
	rabbitMQConfig *config.RabbitMQConfig,
	roomsAPIClient *config.RoomsAPIClient,
	searchService *services.SearchService,
) *RoomsConsumer {
	return &RoomsConsumer{
		rabbitMQConfig: rabbitMQConfig,
		roomsAPIClient: roomsAPIClient,
		searchService:  searchService,
	}
}

// Start inicia el consumer de RabbitMQ
func (c *RoomsConsumer) Start() error {
	// Conectar a RabbitMQ con reintentos
	conn, err := c.rabbitMQConfig.ConnectWithRetry(10, 5*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn

	// Crear channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	c.channel = ch

	// Setup exchange y queue
	if err := c.rabbitMQConfig.SetupExchangeAndQueue(ch); err != nil {
		return err
	}

	// Set QoS para procesar un mensaje a la vez
	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}

	// Consumir mensajes
	msgs, err := ch.Consume(
		c.rabbitMQConfig.QueueName, // queue
		"search-api-consumer",      // consumer tag
		false,                      // auto-ack (false para manual ack)
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	if err != nil {
		return err
	}

	log.Printf("RabbitMQ consumer started, waiting for messages...")

	// Procesar mensajes en goroutines
	go c.processMessages(msgs)

	// Manejar reconexión en caso de cierre de conexión
	go c.handleReconnection()

	return nil
}

// processMessages procesa los mensajes recibidos
func (c *RoomsConsumer) processMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		// Procesar mensaje en goroutine separada
		go c.handleMessage(msg)
	}
}

// handleMessage procesa un mensaje individual
func (c *RoomsConsumer) handleMessage(msg amqp.Delivery) {
	log.Printf("Received message: %s", string(msg.Body))

	var event domain.RoomEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		msg.Nack(false, false) // No requeue si el JSON es inválido
		return
	}

	// Procesar según el tipo de evento
	var err error
	switch event.EventType {
	case "created", "updated":
		err = c.handleCreateOrUpdate(event)
	case "deleted":
		err = c.handleDelete(event)
	default:
		log.Printf("Unknown event type: %s", event.EventType)
		msg.Ack(false) // Ack para no reprocesar
		return
	}

	if err != nil {
		log.Printf("Failed to process event: %v", err)
		// Nack con requeue para reintentar
		msg.Nack(false, true)
		return
	}

	// Ack exitoso
	msg.Ack(false)
	log.Printf("Event processed successfully: %s for room %d", event.EventType, event.RoomID)
}

// handleCreateOrUpdate maneja eventos de creación/actualización
func (c *RoomsConsumer) handleCreateOrUpdate(event domain.RoomEvent) error {
	// Obtener datos completos de la room desde rooms-api
	room, err := c.roomsAPIClient.GetRoomByID(event.RoomID)
	if err != nil {
		return err
	}

	// Indexar en Solr
	if err := c.searchService.IndexRoom(room); err != nil {
		return err
	}

	return nil
}

// handleDelete maneja eventos de eliminación
func (c *RoomsConsumer) handleDelete(event domain.RoomEvent) error {
	// Eliminar del índice de Solr
	if err := c.searchService.DeleteRoom(event.RoomID); err != nil {
		return err
	}

	return nil
}

// handleReconnection maneja la reconexión en caso de cierre de conexión
func (c *RoomsConsumer) handleReconnection() {
	closeChan := make(chan *amqp.Error)
	c.conn.NotifyClose(closeChan)

	closeErr := <-closeChan
	if closeErr != nil {
		log.Printf("RabbitMQ connection closed: %v", closeErr)
		log.Printf("Attempting to reconnect...")

		// Intentar reconectar
		for {
			time.Sleep(5 * time.Second)

			if err := c.Start(); err != nil {
				log.Printf("Failed to reconnect: %v", err)
				continue
			}

			log.Printf("Successfully reconnected to RabbitMQ")
			break
		}
	}
}

// Stop detiene el consumer
func (c *RoomsConsumer) Stop() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	log.Printf("RabbitMQ consumer stopped")
	return nil
}
