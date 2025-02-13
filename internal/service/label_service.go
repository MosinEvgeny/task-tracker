package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/MosinEvgeny/task-tracker/internal/repository"
	"github.com/google/uuid"
)

// LabelService определяет интерфейс для работы с метками.
type LabelService interface {
	CreateLabel(ctx context.Context, name, color string, userID uuid.UUID) (*domain.Label, error)
	GetLabelByID(ctx context.Context, id uuid.UUID) (*domain.Label, error)
	GetAllLabelsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Label, error)
	UpdateLabel(ctx context.Context, id uuid.UUID, name, color string) (*domain.Label, error)
	DeleteLabel(ctx context.Context, id uuid.UUID) error
}

// DefaultLabelService реализует интерфейс LabelService.
type DefaultLabelService struct {
	labelRepo repository.LabelRepository
}

// NewLabelService создает новый экземпляр DefaultLabelService.
func NewLabelService(labelRepo repository.LabelRepository) *DefaultLabelService {
	return &DefaultLabelService{labelRepo: labelRepo}
}

// CreateLabel создает новую метку.
func (s *DefaultLabelService) CreateLabel(ctx context.Context, name, color string, userID uuid.UUID) (*domain.Label, error) {
	if name == "" {
		return nil, fmt.Errorf("необходимо указать название метки")
	}
	if color == "" {
		return nil, fmt.Errorf("необходимо указать цвет метки")
	}
	if userID == uuid.Nil {
		return nil, fmt.Errorf("необходимо указать пользователя")
	}

	hexColorRegex := regexp.MustCompile(`^#([0-9a-fA-F]{3}){1,2}$`)
	if !hexColorRegex.MatchString(color) {
		return nil, fmt.Errorf("неверный формат цвета (HEX)")
	}

	label := &domain.Label{
		ID:     uuid.New(),
		Name:   name,
		Color:  color,
		UserID: userID,
	}

	if err := s.labelRepo.Create(ctx, label); err != nil {
		return nil, fmt.Errorf("ошибка при создании метки: %w", err)
	}

	return label, nil
}

func (s *DefaultLabelService) GetLabelByID(ctx context.Context, id uuid.UUID) (*domain.Label, error) {
	label, err := s.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении метки по ID: %w", err)
	}
	return label, nil
}

func (s *DefaultLabelService) GetAllLabelsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Label, error) {
	labels, err := s.labelRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении меток пользователя: %w", err)
	}
	return labels, nil
}

func (s *DefaultLabelService) UpdateLabel(ctx context.Context, id uuid.UUID, name, color string) (*domain.Label, error) {
	label, err := s.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("метка не найдена")
	}

	label.Name = name
	label.Color = color

	if err := s.labelRepo.Update(ctx, label); err != nil {
		return nil, fmt.Errorf("ошибка при обновлении метки: %w", err)
	}

	return label, nil
}

func (s *DefaultLabelService) DeleteLabel(ctx context.Context, id uuid.UUID) error {
	err := s.labelRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении метки: %w", err)
	}
	return nil
}
