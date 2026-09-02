package repositories

import (
	"time"

	"github.com/karlseguin/ccache/v3"
)

// LocalCacheRepository maneja el caché local (CCache)
type LocalCacheRepository struct {
	cache *ccache.Cache[[]byte]
}

// NewLocalCacheRepository crea un nuevo repositorio de caché local
func NewLocalCacheRepository(cache *ccache.Cache[[]byte]) *LocalCacheRepository {
	return &LocalCacheRepository{
		cache: cache,
	}
}

// Get obtiene un valor del caché local
func (r *LocalCacheRepository) Get(key string) ([]byte, bool) {
	item := r.cache.Get(key)
	if item == nil || item.Expired() {
		return nil, false
	}

	return item.Value(), true
}

// Set almacena un valor en el caché local
func (r *LocalCacheRepository) Set(key string, value []byte, ttl time.Duration) {
	r.cache.Set(key, value, ttl)
}

// Delete elimina un valor del caché local
func (r *LocalCacheRepository) Delete(key string) {
	r.cache.Delete(key)
}

// Clear limpia todo el caché local
func (r *LocalCacheRepository) Clear() {
	r.cache.Clear()
}
