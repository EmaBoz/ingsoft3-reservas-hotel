package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"rooms-api/domain" // 🔴 AJUSTA ESTE IMPORT AL MÓDULO REAL
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type RoomCacheRepository interface {
	Set(ctx context.Context, roomID uint, room domain.RoomResponse) error
	Get(ctx context.Context, roomID uint) (domain.RoomResponse, error)
	Delete(ctx context.Context, roomID uint) error

	SetByNumber(ctx context.Context, number string, room domain.RoomResponse) error
	GetByNumber(ctx context.Context, number string) (domain.RoomResponse, error)
	DeleteByNumber(ctx context.Context, number string) error
}

type roomCacheRepository struct {
	client *memcache.Client
	ttl    time.Duration
}

func NewRoomCacheRepository(host string, port string, ttl time.Duration) RoomCacheRepository {
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))
	return &roomCacheRepository{
		client: client,
		ttl:    ttl,
	}
}

// ---------- Cache por ID ----------

func (r *roomCacheRepository) Set(ctx context.Context, roomID uint, room domain.RoomResponse) error {
	bytes, err := json.Marshal(room)
	if err != nil {
		return fmt.Errorf("error marshalling room to JSON: %w", err)
	}

	key := fmt.Sprintf("room:id:%d", roomID)
	if err := r.client.Set(&memcache.Item{
		Key:        key,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return fmt.Errorf("error setting room in memcached: %w", err)
	}
	return nil
}

func (r *roomCacheRepository) Get(ctx context.Context, roomID uint) (domain.RoomResponse, error) {
	key := fmt.Sprintf("room:id:%d", roomID)
	item, err := r.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return domain.RoomResponse{}, fmt.Errorf("room not found in cache")
		}
		return domain.RoomResponse{}, fmt.Errorf("error getting room from memcached: %w", err)
	}

	var room domain.RoomResponse
	if err := json.Unmarshal(item.Value, &room); err != nil {
		return domain.RoomResponse{}, fmt.Errorf("error unmarshalling room from JSON: %w", err)
	}
	return room, nil
}

func (r *roomCacheRepository) Delete(ctx context.Context, roomID uint) error {
	key := fmt.Sprintf("room:id:%d", roomID)
	if err := r.client.Delete(key); err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("error deleting room from memcached: %w", err)
	}
	return nil
}

// ---------- Cache por Number ----------

func (r *roomCacheRepository) SetByNumber(ctx context.Context, number string, room domain.RoomResponse) error {
	bytes, err := json.Marshal(room)
	if err != nil {
		return fmt.Errorf("error marshalling room to JSON: %w", err)
	}

	key := fmt.Sprintf("room:number:%s", number)
	if err := r.client.Set(&memcache.Item{
		Key:        key,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return fmt.Errorf("error setting room in memcached: %w", err)
	}
	return nil
}

func (r *roomCacheRepository) GetByNumber(ctx context.Context, number string) (domain.RoomResponse, error) {
	key := fmt.Sprintf("room:number:%s", number)
	item, err := r.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return domain.RoomResponse{}, fmt.Errorf("room not found in cache")
		}
		return domain.RoomResponse{}, fmt.Errorf("error getting room from memcached: %w", err)
	}

	var room domain.RoomResponse
	if err := json.Unmarshal(item.Value, &room); err != nil {
		return domain.RoomResponse{}, fmt.Errorf("error unmarshalling room from JSON: %w", err)
	}
	return room, nil
}

func (r *roomCacheRepository) DeleteByNumber(ctx context.Context, number string) error {
	key := fmt.Sprintf("room:number:%s", number)
	if err := r.client.Delete(key); err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("error deleting room from memcached: %w", err)
	}
	return nil
}
