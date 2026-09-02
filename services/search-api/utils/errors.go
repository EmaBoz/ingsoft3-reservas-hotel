package utils

import "errors"

var (
	// Errores de Solr
	ErrSolrConnection    = errors.New("failed to connect to Solr")
	ErrSolrQuery         = errors.New("Solr query failed")
	ErrSolrIndex         = errors.New("failed to index document in Solr")
	ErrSolrDelete        = errors.New("failed to delete document from Solr")

	// Errores de caché
	ErrCacheMiss         = errors.New("cache miss")
	ErrCacheSet          = errors.New("failed to set cache")

	// Errores de validación
	ErrInvalidParameters = errors.New("invalid search parameters")
	ErrInvalidPageSize   = errors.New("page size exceeds maximum allowed")

	// Errores de RabbitMQ
	ErrRabbitMQConnection = errors.New("failed to connect to RabbitMQ")
	ErrRabbitMQPublish    = errors.New("failed to publish message to RabbitMQ")
	ErrRabbitMQConsume    = errors.New("failed to consume message from RabbitMQ")

	// Errores de rooms-api
	ErrRoomNotFound      = errors.New("room not found in rooms-api")
	ErrRoomsAPIRequest   = errors.New("failed to request rooms-api")
)

// AppError representa un error de aplicación con código de estado HTTP
type AppError struct {
	Error      error
	Message    string
	StatusCode int
}

// NewAppError crea un nuevo AppError
func NewAppError(err error, message string, statusCode int) *AppError {
	return &AppError{
		Error:      err,
		Message:    message,
		StatusCode: statusCode,
	}
}
