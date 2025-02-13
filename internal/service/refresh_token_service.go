package service

import (
	"context"
	"fmt"
	"time"

	"os"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/MosinEvgeny/task-tracker/internal/repository"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// RefreshTokenService определяет интерфейс для работы с refresh токенами.
type RefreshTokenService interface {
	CreateRefreshToken(ctx context.Context, userID uuid.UUID) (*domain.RefreshToken, error)
	GetRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) error
	DeleteAllRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) error
}

// DefaultRefreshTokenService реализует интерфейс RefreshTokenService.
type DefaultRefreshTokenService struct {
	refreshTokenRepo repository.RefreshTokenRepository
}

// NewRefreshTokenService создает новый экземпляр DefaultRefreshTokenService.
func NewRefreshTokenService(refreshTokenRepo repository.RefreshTokenRepository) *DefaultRefreshTokenService {
	return &DefaultRefreshTokenService{refreshTokenRepo: refreshTokenRepo}
}

// CreateRefreshToken создает новый refresh токен.
func (s *DefaultRefreshTokenService) CreateRefreshToken(ctx context.Context, userID uuid.UUID) (*domain.RefreshToken, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	RefreshTokenExpireTime := os.Getenv("REFRESH_TOKEN_EXPIRE_TIME")
	duration, err := time.ParseDuration(RefreshTokenExpireTime)
	if err != nil {
		fmt.Println("Error parsing duration:", err)
		return nil, err
	}

	refreshToken := &domain.RefreshToken{
		ID:         uuid.New(),
		UserID:     userID,
		Token:      uuid.New().String(),            // Генерируем случайный токен
		ExpiryDate: time.Now().UTC().Add(duration), // Срок действия 7 дней
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("ошибка при создании refresh токена: %w", err)
	}

	return refreshToken, nil
}

// GetRefreshToken получает refresh токен по токену.
func (s *DefaultRefreshTokenService) GetRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	refreshToken, err := s.refreshTokenRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении refresh токена: %w", err)
	}

	return refreshToken, nil
}

// DeleteRefreshToken удаляет refresh токен.
func (s *DefaultRefreshTokenService) DeleteRefreshToken(ctx context.Context, id uuid.UUID) error {
	err := s.refreshTokenRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении refresh токена: %w", err)
	}

	return nil
}

// DeleteAllRefreshTokensByUserID удаляет все refresh токены пользователя.
func (s *DefaultRefreshTokenService) DeleteAllRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	err := s.refreshTokenRepo.DeleteAllByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении всех refresh токенов пользователя: %w", err)
	}

	return nil
}
