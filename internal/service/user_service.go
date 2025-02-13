package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/MosinEvgeny/task-tracker/internal/repository"
	"github.com/google/uuid"
)

// UserService определяет интерфейс для работы с пользователями.
type UserService interface {
	CreateUser(ctx context.Context, username, email, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, username, email string) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// DefaultUserService реализует интерфейс UserService.
type DefaultUserService struct {
	userRepo repository.UserRepository
}

// NewUserService создает новый экземпляр DefaultUserService.
func NewUserService(userRepo repository.UserRepository) *DefaultUserService {
	return &DefaultUserService{userRepo: userRepo}
}

func (s *DefaultUserService) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, fmt.Errorf("необходимо заполнить все поля")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, fmt.Errorf("неверный формат email")
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("пользователь с таким email уже существует")
	}

	user := &domain.User{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
	}

	if err := user.HashPassword(password); err != nil {
		return nil, fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	return user, nil
}

func (s *DefaultUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя по ID: %w", err)
	}
	return user, nil
}

func (s *DefaultUserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя по email: %w", err)
	}
	return user, nil
}

func (s *DefaultUserService) UpdateUser(ctx context.Context, id uuid.UUID, username, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден")
	}

	user.Username = username
	user.Email = email

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	return user, nil
}

func (s *DefaultUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}
	return nil
}
