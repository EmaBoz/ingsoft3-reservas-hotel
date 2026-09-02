package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"users-api/domain"

	"github.com/bradfitz/gomemcache/memcache"
)

type UserCacheRepository interface {
	Set(ctx context.Context, userID uint, user domain.UserResponseDTO) error
	Get(ctx context.Context, userID uint) (domain.UserResponseDTO, error)
	Delete(ctx context.Context, userID uint) error
	SetByUsername(ctx context.Context, username string, user domain.UserResponseDTO) error
	GetByUsername(ctx context.Context, username string) (domain.UserResponseDTO, error)
	DeleteByUsername(ctx context.Context, username string) error
}

type userCacheRepository struct {
	client *memcache.Client
	ttl    time.Duration
}

func NewUserCacheRepository(host string, port string, ttl time.Duration) UserCacheRepository {
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))
	return &userCacheRepository{
		client: client,
		ttl:    ttl,
	}
}

// Set guarda un usuario en cache por su ID
func (r *userCacheRepository) Set(ctx context.Context, userID uint, user domain.UserResponseDTO) error {
	bytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshalling user to JSON: %w", err)
	}

	key := fmt.Sprintf("user:id:%d", userID)
	if err := r.client.Set(&memcache.Item{
		Key:        key,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return fmt.Errorf("error setting user in memcached: %w", err)
	}
	return nil
}

// Get obtiene un usuario desde cache por su ID
func (r *userCacheRepository) Get(ctx context.Context, userID uint) (domain.UserResponseDTO, error) {
	key := fmt.Sprintf("user:id:%d", userID)
	item, err := r.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return domain.UserResponseDTO{}, fmt.Errorf("user not found in cache")
		}
		return domain.UserResponseDTO{}, fmt.Errorf("error getting user from memcached: %w", err)
	}

	var user domain.UserResponseDTO
	if err := json.Unmarshal(item.Value, &user); err != nil {
		return domain.UserResponseDTO{}, fmt.Errorf("error unmarshalling user from JSON: %w", err)
	}
	return user, nil
}

// Delete elimina un usuario del cache por su ID
func (r *userCacheRepository) Delete(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user:id:%d", userID)
	if err := r.client.Delete(key); err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("error deleting user from memcached: %w", err)
	}
	return nil
}

// SetByUsername guarda un usuario en cache por su username
func (r *userCacheRepository) SetByUsername(ctx context.Context, username string, user domain.UserResponseDTO) error {
	bytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshalling user to JSON: %w", err)
	}

	key := fmt.Sprintf("user:username:%s", username)
	if err := r.client.Set(&memcache.Item{
		Key:        key,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return fmt.Errorf("error setting user in memcached: %w", err)
	}
	return nil
}

// GetByUsername obtiene un usuario desde cache por su username
func (r *userCacheRepository) GetByUsername(ctx context.Context, username string) (domain.UserResponseDTO, error) {
	key := fmt.Sprintf("user:username:%s", username)
	item, err := r.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return domain.UserResponseDTO{}, fmt.Errorf("user not found in cache")
		}
		return domain.UserResponseDTO{}, fmt.Errorf("error getting user from memcached: %w", err)
	}

	var user domain.UserResponseDTO
	if err := json.Unmarshal(item.Value, &user); err != nil {
		return domain.UserResponseDTO{}, fmt.Errorf("error unmarshalling user from JSON: %w", err)
	}
	return user, nil
}

// DeleteByUsername elimina un usuario del cache por su username
func (r *userCacheRepository) DeleteByUsername(ctx context.Context, username string) error {
	key := fmt.Sprintf("user:username:%s", username)
	if err := r.client.Delete(key); err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("error deleting user from memcached: %w", err)
	}
	return nil
}
