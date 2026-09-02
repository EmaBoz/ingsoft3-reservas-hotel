package services

import (
	"context"
	"fmt"
	"os"
	"time"
	"users-api/domain"
	"users-api/repositories"
	"users-api/utils"
)

type UserService interface {
	GetAllUsers() ([]domain.UserResponseDTO, error)
	CreateUser(ctx context.Context, dto domain.CreateUserDTO) (domain.UserResponseDTO, error)
	GetUserByID(id uint) (domain.UserResponseDTO, error)
	Login(dto domain.LoginDTO) (domain.LoginResponseDTO, error)
}

type userService struct {
	repository repositories.UserRepository
	cache      repositories.UserCacheRepository
}

func NewUserService(repository repositories.UserRepository, cache repositories.UserCacheRepository) UserService {

	return &userService{
		repository: repository,
		cache:      cache,
	}
}

func (s *userService) GetAllUsers() ([]domain.UserResponseDTO, error) {
	return s.repository.GetAllUsers()
}

func (s *userService) CreateUser(ctx context.Context, dto domain.CreateUserDTO) (domain.UserResponseDTO, error) {
	// Hashear la contraseña antes de guardarla
	hashed, err := utils.HashPassword(dto.Password)
	if err != nil {
		return domain.UserResponseDTO{}, fmt.Errorf("failed to hash password: %w", err)
	}
	dto.Password = hashed

	created, err := s.repository.Create(ctx, dto)
	if err != nil {
		return domain.UserResponseDTO{}, fmt.Errorf("failed to create user: %w", err)
	}
	err = s.cache.Set(ctx, created.ID, created)
	if err != nil {
		// Loguear el error pero no interrumpir el flujo
		fmt.Printf("warning: failed to cache user after creation: %v\n", err)
	}

	return created, nil
}

func (s *userService) GetUserByID(id uint) (domain.UserResponseDTO, error) {
	ctx := context.Background()

	// 1) Intento de caché
	if cached, err := s.cache.Get(ctx, id); err == nil {
		return cached, nil
	}

	// 2) Fallback a DB
	user, err := s.repository.GetByID(id)
	if err != nil {
		return domain.UserResponseDTO{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// 3) Rehidratar caché (best-effort)
	_ = s.cache.Set(ctx, id, user)

	return user, nil
}

func (s *userService) Login(dto domain.LoginDTO) (domain.LoginResponseDTO, error) {
	user, err := s.repository.GetByUsernameOrEmail(dto.UsernameOrEmail)
	if err != nil {
		return domain.LoginResponseDTO{}, fmt.Errorf("invalid credentials")
	}

	if err := utils.CheckPassword(user.Password, dto.Password); err != nil {
		return domain.LoginResponseDTO{}, fmt.Errorf("invalid credentials")
	}

	secret := os.Getenv("JWT_SECRET")
	token, err := utils.GenerateToken(user, secret)
	if err != nil {
		return domain.LoginResponseDTO{}, fmt.Errorf("failed to generate token: %w", err)
	}

	// decí cuánto dura el token: en jwt.go pusimos 24h
	expiresAt := time.Now().Add(24 * time.Hour)

	// Guardar token en DB
	if err := s.repository.UpdateToken(user.ID, token, &expiresAt); err != nil {
		// Si falla persistir el token, podés elegir si igual retornás el token o no.
		return domain.LoginResponseDTO{}, fmt.Errorf("failed to persist token: %w", err)
	}

	resp := domain.LoginResponseDTO{
		Token: token,
		User: domain.UserResponseDTO{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
		},
	}
	return resp, nil
}
