package repository

import (
	"context"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/google/uuid"
)

// TaskRepository определяет интерфейс для работы с задачами в базе данных.
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Task, error) // Получение всех задач пользователя
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}
