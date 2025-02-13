package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/MosinEvgeny/task-tracker/internal/repository"
	"github.com/google/uuid"
)

// TaskService определяет интерфейс для работы с задачами.
type TaskService interface {
	CreateTask(ctx context.Context, title, description string, dueDate time.Time, userID uuid.UUID) (*domain.Task, error)
	GetTaskByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetAllTasksByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Task, error)
	UpdateTask(ctx context.Context, id uuid.UUID, title, description string, dueDate time.Time) (*domain.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
}

// DefaultTaskService реализует интерфейс TaskService.
type DefaultTaskService struct {
	taskRepo repository.TaskRepository
}

// NewTaskService создает новый экземпляр DefaultTaskService.
func NewTaskService(taskRepo repository.TaskRepository) *DefaultTaskService {
	return &DefaultTaskService{taskRepo: taskRepo}
}

// CreateTask создает новую задачу.
func (s *DefaultTaskService) CreateTask(ctx context.Context, title, description string, dueDate time.Time, userID uuid.UUID) (*domain.Task, error) {
	if title == "" {
		return nil, fmt.Errorf("необходимо указать название задачи")
	}
	if userID == uuid.Nil {
		return nil, fmt.Errorf("необходимо указать пользователя")
	}

	task := &domain.Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		UserID:      userID,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("ошибка при создании задачи: %w", err)
	}

	return task, nil
}

func (s *DefaultTaskService) GetTaskByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задачи по ID: %w", err)
	}
	return task, nil
}

func (s *DefaultTaskService) GetAllTasksByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач пользователя: %w", err)
	}
	return tasks, nil
}

func (s *DefaultTaskService) UpdateTask(ctx context.Context, id uuid.UUID, title, description string, dueDate time.Time) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}

	task.Title = title
	task.Description = description
	task.DueDate = dueDate

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("ошибка при обновлении задачи: %w", err)
	}

	return task, nil
}

func (s *DefaultTaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	err := s.taskRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}
	return nil
}
