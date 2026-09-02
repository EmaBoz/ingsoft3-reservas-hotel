package services

import (
	"context"
	"fmt"
	"net/url"
	"rooms-api/config"
	"rooms-api/domain"
	"rooms-api/events"
	"rooms-api/repositories"
	"strconv"
)

type RoomService struct {
	roomRepo        *repositories.RoomRepository
	publisher       *events.EventPublisher
	cache           repositories.RoomCacheRepository
	searchAPIClient *config.SearchAPIClient // 🔹 NUEVO
}

func NewRoomService(
	roomRepo *repositories.RoomRepository,
	publisher *events.EventPublisher,
	cache repositories.RoomCacheRepository,
	searchAPIClient *config.SearchAPIClient, // 🔹 NUEVO
) *RoomService {
	return &RoomService{
		roomRepo:        roomRepo,
		publisher:       publisher,
		cache:           cache,
		searchAPIClient: searchAPIClient, // 🔹 NUEVO
	}
}

func (s *RoomService) CreateRoom(ctx context.Context, req domain.CreateRoomRequest) (*domain.RoomResponse, error) {
	room := &domain.Room{
		Number:      req.Number,
		Type:        req.Type,
		Status:      domain.RoomStatusAvailable,
		Price:       req.Price,
		Description: req.Description,
		Capacity:    req.Capacity,
		Floor:       req.Floor,
		HasWifi:     req.HasWifi,
		HasAC:       req.HasAC,
		HasTV:       req.HasTV,
		HasMinibar:  req.HasMinibar,
	}

	err := s.roomRepo.Create(ctx, room)
	if err != nil {
		return nil, err
	}

	err = s.cache.Set(ctx, room.ID, *s.roomToResponse(room))
	if err != nil {
		return nil, err
	}

	// Publicar evento de creación
	if s.publisher != nil {
		go s.publisher.PublishRoomCreated(room)
	}

	return s.roomToResponse(room), nil
}

func (s *RoomService) GetRoomByID(ctx context.Context, id uint) (*domain.RoomResponse, error) {
	cached, err := s.cache.Get(ctx, id)
	if err == nil {
		return &cached, nil
	}

	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Try to update the cache but don't fail the request if caching fails
	_ = s.cache.Set(ctx, room.ID, *s.roomToResponse(room))

	return s.roomToResponse(room), nil
}

func (s *RoomService) GetRoomByNumber(ctx context.Context, number string) (*domain.RoomResponse, error) {
	room, err := s.roomRepo.GetByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	return s.roomToResponse(room), nil
}

// ✅ MODIFICADO: Ahora usa Search API en lugar de consulta directa a MySQL
func (s *RoomService) GetRooms(ctx context.Context, filter domain.RoomFilter, page, limit int) (*domain.RoomListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 🔹 DELEGAR A SEARCH API para búsquedas con filtros
	if s.hasComplexFilters(filter) {
		// Usar Search API (Solr) para búsquedas complejas
		return s.searchAPIClient.SearchRooms(filter, page, limit)
	}

	// Para búsquedas simples sin filtros, usar MySQL directamente (más rápido)
	rooms, total, err := s.roomRepo.GetAll(ctx, filter, page, limit)
	if err != nil {
		return nil, err
	}

	roomResponses := make([]domain.RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = *s.roomToResponse(&room)
	}

	return &domain.RoomListResponse{
		Rooms: roomResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// hasComplexFilters determina si debe usar Search API (Solr) o MySQL directo
func (s *RoomService) hasComplexFilters(filter domain.RoomFilter) bool {
	// Usar Search API (Solr) cuando hay:
	// - Búsqueda de texto libre (no implementado en MySQL)
	// - Múltiples filtros combinados (más eficiente en Solr)
	// - Filtros de amenities (has_wifi, has_ac, etc.)

	filterCount := 0

	if filter.Type != nil {
		filterCount++
	}
	if filter.Status != nil {
		filterCount++
	}
	if filter.Floor != nil {
		filterCount++
	}
	if filter.MinPrice != nil || filter.MaxPrice != nil {
		filterCount++
	}
	if filter.HasWifi != nil || filter.HasAC != nil || filter.HasTV != nil || filter.HasMinibar != nil {
		return true // Siempre usar Solr para filtros de amenities
	}

	// Si hay más de 2 filtros, usar Solr para mejor performance
	return filterCount > 2
}

func (s *RoomService) UpdateRoom(ctx context.Context, id uint, req domain.UpdateRoomRequest) (*domain.RoomResponse, error) {
	// Obtener estado anterior para detectar cambios de status
	oldRoom, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	oldStatus := oldRoom.Status

	// Actualizar
	room, err := s.roomRepo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Invalidar caché
	_ = s.cache.Delete(ctx, room.ID)

	// Publicar evento de actualización
	if s.publisher != nil {
		go s.publisher.PublishRoomUpdated(room)

		// Si cambió el status, publicar evento específico
		if oldStatus != room.Status {
			go s.publisher.PublishRoomStatusChanged(room.ID, string(oldStatus), string(room.Status))
		}
	}

	return s.roomToResponse(room), nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, id uint) error {
	// Obtener ID antes de eliminar
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.roomRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidar caché
	_ = s.cache.Delete(ctx, room.ID)

	// Publicar evento de eliminación
	if s.publisher != nil {
		go s.publisher.PublishRoomDeleted(room.ID)
	}

	return nil
}

func (s *RoomService) UpdateRoomStatus(ctx context.Context, id uint, status domain.RoomStatus) (*domain.RoomResponse, error) {
	updateReq := domain.UpdateRoomRequest{
		Status: &status,
	}
	return s.UpdateRoom(ctx, id, updateReq)
}

// ✅ MODIFICADO: GetAvailableRooms ahora usa Search API
func (s *RoomService) GetAvailableRooms(ctx context.Context, filter domain.RoomFilter, page, limit int) (*domain.RoomListResponse, error) {
	availableStatus := domain.RoomStatusAvailable
	filter.Status = &availableStatus

	// Siempre usar Search API para habitaciones disponibles (consulta frecuente)
	return s.searchAPIClient.SearchRooms(filter, page, limit)
}

func (s *RoomService) GetRoomsByType(ctx context.Context, roomType domain.RoomType, page, limit int) (*domain.RoomListResponse, error) {
	filter := domain.RoomFilter{
		Type: &roomType,
	}
	return s.GetRooms(ctx, filter, page, limit)
}

func (s *RoomService) GetRoomsByFloor(ctx context.Context, floor int, page, limit int) (*domain.RoomListResponse, error) {
	filter := domain.RoomFilter{
		Floor: &floor,
	}
	return s.GetRooms(ctx, filter, page, limit)
}

func (s *RoomService) SearchRooms(ctx context.Context, query string, page, limit int) (*domain.RoomListResponse, error) {
	// Búsqueda de texto libre siempre usa Search API
	filter := domain.RoomFilter{}
	return s.searchAPIClient.SearchRooms(filter, page, limit)
}

func (s *RoomService) roomToResponse(room *domain.Room) *domain.RoomResponse {
	return &domain.RoomResponse{
		ID:          room.ID,
		Number:      room.Number,
		Type:        room.Type,
		Status:      room.Status,
		Price:       room.Price,
		Description: room.Description,
		Capacity:    room.Capacity,
		Floor:       room.Floor,
		HasWifi:     room.HasWifi,
		HasAC:       room.HasAC,
		HasTV:       room.HasTV,
		HasMinibar:  room.HasMinibar,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
	}
}

// GetRoomsViaSearch busca habitaciones usando Search API (Solr)
func (s *RoomService) GetRoomsViaSearch(ctx context.Context, filter domain.RoomFilter, page, limit int) (*domain.RoomListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Construir parámetros para Search API
	params := url.Values{}

	if filter.Type != nil {
		params.Add("type", string(*filter.Type))
	}
	if filter.Status != nil {
		params.Add("status", string(*filter.Status))
	}
	if filter.Floor != nil {
		params.Add("floor", strconv.Itoa(*filter.Floor))
	}
	if filter.MinPrice != nil {
		params.Add("min_price", fmt.Sprintf("%.2f", *filter.MinPrice))
	}
	if filter.MaxPrice != nil {
		params.Add("max_price", fmt.Sprintf("%.2f", *filter.MaxPrice))
	}
	if filter.HasWifi != nil {
		params.Add("has_wifi", strconv.FormatBool(*filter.HasWifi))
	}
	if filter.HasAC != nil {
		params.Add("has_ac", strconv.FormatBool(*filter.HasAC))
	}
	if filter.HasTV != nil {
		params.Add("has_tv", strconv.FormatBool(*filter.HasTV))
	}
	if filter.HasMinibar != nil {
		params.Add("has_minibar", strconv.FormatBool(*filter.HasMinibar))
	}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	// Llamar a Search API
	respData, err := s.searchAPIClient.SearchRooms(filter, page, limit)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
