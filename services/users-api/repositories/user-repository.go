package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"
	"users-api/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUsers() ([]domain.UserResponseDTO, error)
	Create(ctx context.Context, user domain.CreateUserDTO) (domain.UserResponseDTO, error)
	GetByID(id uint) (domain.UserResponseDTO, error)
	UpdateToken(userID uint, token string, expiresAt *time.Time) error
	// nuevo: obtener usuario completo (incluye password) por username o email
	GetByUsernameOrEmail(identifier string) (domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers() ([]domain.UserResponseDTO, error) {
	var users []domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	var userResponses []domain.UserResponseDTO
	for _, u := range users {
		userResponses = append(userResponses, domain.UserResponseDTO{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
		})
	}
	return userResponses, nil
}

func (r *userRepository) Create(ctx context.Context, user domain.CreateUserDTO) (domain.UserResponseDTO, error) {
	newUser := domain.User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password, // ya viene hasheada por el service
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      domain.RoleNormal, // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := r.db.Create(&newUser).Error; err != nil {
		return domain.UserResponseDTO{}, fmt.Errorf("failed to create user: %w", err)
	}

	response := domain.UserResponseDTO{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Role:      newUser.Role,
	}

	return response, nil
}

func (r *userRepository) GetByID(id uint) (domain.UserResponseDTO, error) {
	var u domain.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.UserResponseDTO{}, gorm.ErrRecordNotFound
		}
		return domain.UserResponseDTO{}, fmt.Errorf("failed to get user by id: %w", err)
	}
	return domain.UserResponseDTO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
	}, nil
}

// Nuevo: GetByUsernameOrEmail -> retorna domain.User (incluye password)
func (r *userRepository) GetByUsernameOrEmail(identifier string) (domain.User, error) {
	var u domain.User
	if err := r.db.Where("username = ? OR email = ?", identifier, identifier).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, gorm.ErrRecordNotFound
		}
		return domain.User{}, fmt.Errorf("failed to get user by username or email: %w", err)
	}
	return u, nil
}

func (r *userRepository) UpdateToken(userID uint, token string, expiresAt *time.Time) error {
	updates := map[string]interface{}{
		"token": token,
	}
	if expiresAt != nil {
		updates["token_expires_at"] = *expiresAt
	} else {
		updates["token_expires_at"] = nil
	}

	if err := r.db.Model(&domain.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}
	return nil
}
