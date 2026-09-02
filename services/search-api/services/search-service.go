package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"search-api/config"
	"search-api/domain"
	"search-api/repositories"
)

// SearchService maneja la lógica de negocio de búsqueda
type SearchService struct {
	solrRepo        *repositories.SolrRepository
	localCache      *repositories.LocalCacheRepository
	distributedCache *repositories.DistributedCacheRepository
	cacheConfig     *config.CacheConfig
}

// NewSearchService crea un nuevo servicio de búsqueda
func NewSearchService(
	solrRepo *repositories.SolrRepository,
	localCache *repositories.LocalCacheRepository,
	distributedCache *repositories.DistributedCacheRepository,
	cacheConfig *config.CacheConfig,
) *SearchService {
	return &SearchService{
		solrRepo:        solrRepo,
		localCache:      localCache,
		distributedCache: distributedCache,
		cacheConfig:     cacheConfig,
	}
}

// SearchRooms busca habitaciones con doble caché
func (s *SearchService) SearchRooms(req *domain.SearchRoomsRequest) (*domain.SearchRoomsResponse, error) {
	// Validar parámetros
	if err := s.validateSearchRequest(req); err != nil {
		return nil, err
	}

	// Normalizar parámetros (defaults)
	s.normalizeSearchRequest(req)

	// Generar clave de caché
	cacheKey := s.generateCacheKey(req)

	// 1. Intentar caché local (CCache)
	if cachedData, found := s.localCache.Get(cacheKey); found {
		log.Printf("Cache HIT (local): %s", cacheKey)
		var response domain.SearchRoomsResponse
		if err := json.Unmarshal(cachedData, &response); err == nil {
			return &response, nil
		}
		log.Printf("Failed to unmarshal local cache data, continuing...")
	}

	// 2. Intentar caché distribuido (Memcached)
	if cachedData, err := s.distributedCache.Get(cacheKey); err == nil && cachedData != nil {
		log.Printf("Cache HIT (distributed): %s", cacheKey)
		var response domain.SearchRoomsResponse
		if err := json.Unmarshal(cachedData, &response); err == nil {
			// Guardar en caché local para próximas consultas
			s.localCache.Set(cacheKey, cachedData, s.cacheConfig.LocalTTL)
			return &response, nil
		}
		log.Printf("Failed to unmarshal distributed cache data, continuing...")
	}

	// 3. Cache MISS - Consultar Solr
	log.Printf("Cache MISS: %s - Querying Solr", cacheKey)
	response, err := s.solrRepo.Search(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search in Solr: %w", err)
	}

	// 4. Guardar en cachés
	if responseData, err := json.Marshal(response); err == nil {
		// Guardar en Memcached primero (TTL más largo)
		if err := s.distributedCache.Set(cacheKey, responseData, s.cacheConfig.DistributedTTL); err != nil {
			log.Printf("Failed to set distributed cache: %v", err)
		}

		// Luego en caché local (TTL más corto)
		s.localCache.Set(cacheKey, responseData, s.cacheConfig.LocalTTL)
	}

	return response, nil
}

// IndexRoom indexa una habitación en Solr
func (s *SearchService) IndexRoom(room *domain.Room) error {
	doc := room.ToSolrDocument()

	if err := s.solrRepo.IndexDocument(doc); err != nil {
		return fmt.Errorf("failed to index room: %w", err)
	}

	// Invalidar caché (limpiar todo para simplificar)
	s.invalidateCache()

	log.Printf("Room %d indexed successfully in Solr", room.ID)
	return nil
}

// DeleteRoom elimina una habitación del índice de Solr
func (s *SearchService) DeleteRoom(roomID uint) error {
	// Convertir ID a string
	id := fmt.Sprintf("%d", roomID)

	if err := s.solrRepo.DeleteDocument(id); err != nil {
		return fmt.Errorf("failed to delete room from Solr: %w", err)
	}

	// Invalidar caché
	s.invalidateCache()

	log.Printf("Room %d deleted successfully from Solr", roomID)
	return nil
}

// validateSearchRequest valida los parámetros de búsqueda
func (s *SearchService) validateSearchRequest(req *domain.SearchRoomsRequest) error {
	if req.Limit > 50 {
		return fmt.Errorf("limit exceeds maximum allowed (50)")
	}

	if req.Limit < 0 {
		return fmt.Errorf("limit must be positive")
	}

	if req.Page < 0 {
		return fmt.Errorf("page must be positive")
	}

	if req.MinPrice != nil && req.MaxPrice != nil && *req.MinPrice > *req.MaxPrice {
		return fmt.Errorf("min_price cannot be greater than max_price")
	}

	return nil
}

// normalizeSearchRequest establece valores por defecto
func (s *SearchService) normalizeSearchRequest(req *domain.SearchRoomsRequest) {
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}
}

// generateCacheKey genera una clave de caché única basada en los parámetros de búsqueda
func (s *SearchService) generateCacheKey(req *domain.SearchRoomsRequest) string {
	// Serializar request a JSON para tener una representación consistente
	data, err := json.Marshal(req)
	if err != nil {
		// Fallback simple si falla la serialización
		return fmt.Sprintf("search:%s:%s:%d:%d", req.Q, req.Type, req.Page, req.Limit)
	}

	// Hash SHA-256 para tener claves cortas y consistentes
	hash := sha256.Sum256(data)
	return fmt.Sprintf("search:%x", hash[:16]) // Usar solo los primeros 16 bytes
}

// invalidateCache invalida todo el caché (estrategia simple)
func (s *SearchService) invalidateCache() {
	// En producción podrías implementar una estrategia más granular
	// Por ahora, simplemente limpiamos el caché local
	// Memcached se puede dejar que expire naturalmente
	s.localCache.Clear()
	log.Printf("Cache invalidated")
}

// HealthCheck verifica la salud del servicio
func (s *SearchService) HealthCheck() error {
	// Verificar Solr
	if err := s.solrRepo.HealthCheck(); err != nil {
		return fmt.Errorf("Solr health check failed: %w", err)
	}

	// Verificar Memcached
	if err := s.distributedCache.HealthCheck(); err != nil {
		return fmt.Errorf("Memcached health check failed: %w", err)
	}

	return nil
}
