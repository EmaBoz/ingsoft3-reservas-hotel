package repositories

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// DistributedCacheRepository maneja el caché distribuido (Memcached)
type DistributedCacheRepository struct {
	client *memcache.Client
}

// NewDistributedCacheRepository crea un nuevo repositorio de caché distribuido
func NewDistributedCacheRepository(client *memcache.Client) *DistributedCacheRepository {
	return &DistributedCacheRepository{
		client: client,
	}
}

// Get obtiene un valor del caché distribuido
func (r *DistributedCacheRepository) Get(key string) ([]byte, error) {
	item, err := r.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil // No es un error, simplemente no existe
		}
		return nil, fmt.Errorf("failed to get from Memcached: %w", err)
	}

	return item.Value, nil
}

// Set almacena un valor en el caché distribuido
func (r *DistributedCacheRepository) Set(key string, value []byte, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(ttl.Seconds()),
	}

	if err := r.client.Set(item); err != nil {
		return fmt.Errorf("failed to set in Memcached: %w", err)
	}

	return nil
}

// Delete elimina un valor del caché distribuido
func (r *DistributedCacheRepository) Delete(key string) error {
	if err := r.client.Delete(key); err != nil {
		if err == memcache.ErrCacheMiss {
			return nil // No es un error si no existe
		}
		return fmt.Errorf("failed to delete from Memcached: %w", err)
	}

	return nil
}

// HealthCheck verifica la conexión con Memcached
func (r *DistributedCacheRepository) HealthCheck() error {
	// Intentar hacer ping con un Get de una clave que no existe
	_, err := r.client.Get("health-check-ping")
	if err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("Memcached health check failed: %w", err)
	}
	return nil
}
