package repository

import (
	"context"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/google/uuid"
)

// LabelRepository определяет интерфейс для работы с метками в базе данных.
type LabelRepository interface {
	Create(ctx context.Context, label *domain.Label) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Label, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Label, error)
	Update(ctx context.Context, label *domain.Label) error
	Delete(ctx context.Context, id uuid.UUID) error
}
