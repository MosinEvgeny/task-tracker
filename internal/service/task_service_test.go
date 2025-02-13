package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository - это mock для TaskRepository.
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	args := m.Called(ctx, id)
	task, ok := args.Get(0).(*domain.Task)
	if !ok {
		return nil, args.Error(1)
	}
	return task, args.Error(1)
}

func (m *MockTaskRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Task, error) {
	args := m.Called(ctx, userID)
	tasks, ok := args.Get(0).([]*domain.Task)
	if !ok {
		return nil, args.Error(1)
	}
	return tasks, args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	title := "Test Task"
	description := "Test Description"
	dueDate := time.Now()
	userID := uuid.New()

	// Настройка mock-репозитория
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil)

	// 2. Act
	task, err := taskService.CreateTask(ctx, title, description, dueDate, userID)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, title, task.Title)
	assert.Equal(t, description, task.Description)
	assert.Equal(t, dueDate, task.DueDate)
	assert.Equal(t, userID, task.UserID)

	mockRepo.AssertExpectations(t)
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	title := ""
	description := "Test Description"
	dueDate := time.Now()
	userID := uuid.New()

	// 2. Act
	task, err := taskService.CreateTask(ctx, title, description, dueDate, userID)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.EqualError(t, err, "необходимо указать название задачи")

	mockRepo.AssertExpectations(t)
}

func TestGetTaskByID(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	taskID := uuid.New()
	expectedTask := &domain.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     time.Now(),
		UserID:      uuid.New(),
	}

	// Настройка mock-репозитория
	mockRepo.On("GetByID", mock.Anything, taskID).Return(expectedTask, nil)

	// 2. Act
	task, err := taskService.GetTaskByID(ctx, taskID)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTask, task)

	mockRepo.AssertExpectations(t)
}

func TestGetTaskByID_NotFound(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	taskID := uuid.New()

	// Настройка mock-репозитория
	mockRepo.On("GetByID", mock.Anything, taskID).Return(nil, errors.New("task not found"))

	// 2. Act
	task, err := taskService.GetTaskByID(ctx, taskID)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.EqualError(t, err, "ошибка при получении задачи по ID: task not found")

	mockRepo.AssertExpectations(t)
}

func TestUpdateTask(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	taskID := uuid.New()
	initialTask := &domain.Task{
		ID:          taskID,
		Title:       "Old Title",
		Description: "Old Description",
		DueDate:     time.Now(),
		UserID:      uuid.New(),
	}
	updatedTitle := "New Title"
	updatedDescription := "New Description"
	updatedDueDate := time.Now().Add(time.Hour * 24)

	// Настройка mock-репозитория
	mockRepo.On("GetByID", mock.Anything, taskID).Return(initialTask, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(task *domain.Task) bool {
		return task.ID == taskID && task.Title == updatedTitle && task.Description == updatedDescription && task.DueDate == updatedDueDate
	})).Return(nil)

	// 2. Act
	task, err := taskService.UpdateTask(ctx, taskID, updatedTitle, updatedDescription, updatedDueDate)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, taskID, task.ID)
	assert.Equal(t, updatedTitle, task.Title)
	assert.Equal(t, updatedDescription, task.Description)
	assert.Equal(t, updatedDueDate, task.DueDate)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTask_NotFound(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	taskID := uuid.New()
	updatedTitle := "New Title"
	updatedDescription := "New Description"
	updatedDueDate := time.Now().Add(time.Hour * 24)

	// Настройка mock-репозитория
	mockRepo.On("GetByID", mock.Anything, taskID).Return(nil, errors.New("task not found"))

	// 2. Act
	task, err := taskService.UpdateTask(ctx, taskID, updatedTitle, updatedDescription, updatedDueDate)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.EqualError(t, err, "задача не найдена")

	mockRepo.AssertExpectations(t)
}

func TestDeleteTask(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	taskID := uuid.New()

	// Настройка mock-репозитория
	mockRepo.On("Delete", mock.Anything, taskID).Return(nil)

	// 2. Act
	err := taskService.DeleteTask(ctx, taskID)

	// 3. Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteTask_Error(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	taskID := uuid.New()
	expectedError := errors.New("delete error")

	// Настройка mock-репозитория
	mockRepo.On("Delete", mock.Anything, taskID).Return(expectedError)

	// 2. Act
	err := taskService.DeleteTask(ctx, taskID)

	// 3. Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "ошибка при удалении задачи: delete error")

	mockRepo.AssertExpectations(t)
}

func TestGetAllTasksByUserID(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	expectedTasks := []*domain.Task{
		{ID: uuid.New(), Title: "Task 1", Description: "Description 1", DueDate: time.Now(), UserID: userID},
		{ID: uuid.New(), Title: "Task 2", Description: "Description 2", DueDate: time.Now(), UserID: userID},
	}

	// Настройка mock-репозитория
	mockRepo.On("GetAllByUserID", mock.Anything, userID).Return(expectedTasks, nil)

	// 2. Act
	tasks, err := taskService.GetAllTasksByUserID(ctx, userID)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, tasks)

	mockRepo.AssertExpectations(t)
}

func TestGetAllTasksByUserID_Error(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	expectedError := errors.New("get all error")

	// Настройка mock-репозитория
	mockRepo.On("GetAllByUserID", mock.Anything, userID).Return(nil, expectedError)

	// 2. Act
	tasks, err := taskService.GetAllTasksByUserID(ctx, userID)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, tasks)
	assert.EqualError(t, err, "ошибка при получении задач пользователя: get all error")

	mockRepo.AssertExpectations(t)
}
