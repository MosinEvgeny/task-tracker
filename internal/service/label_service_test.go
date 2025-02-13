package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MosinEvgeny/task-tracker/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLabelRepository - это mock для LabelRepository.
type MockLabelRepository struct {
	mock.Mock
}

func (m *MockLabelRepository) Create(ctx context.Context, label *domain.Label) error {
	args := m.Called(ctx, label)
	return args.Error(0)
}

func (m *MockLabelRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Label, error) {
	args := m.Called(ctx, id)
	label, ok := args.Get(0).(*domain.Label)
	if !ok {
		return nil, args.Error(1)
	}
	return label, args.Error(1)
}

func (m *MockLabelRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Label, error) {
	args := m.Called(ctx, userID)
	labels, ok := args.Get(0).([]*domain.Label)
	if !ok {
		return nil, args.Error(1)
	}
	return labels, args.Error(1)
}

func (m *MockLabelRepository) Update(ctx context.Context, label *domain.Label) error {
	args := m.Called(ctx, label)
	return args.Error(0)
}

func (m *MockLabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateLabel(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	name := "Test Label"
	color := "#FFFFFF"
	userID := uuid.New()

	// Настройка mock-репозитория
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Label")).Return(nil)

	// 2. Act
	label, err := labelService.CreateLabel(ctx, name, color, userID)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotNil(t, label)
	assert.Equal(t, name, label.Name)
	assert.Equal(t, color, label.Color)
	assert.Equal(t, userID, label.UserID)

	mockRepo.AssertExpectations(t)
}

func TestCreateLabel_EmptyName(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	name := ""
	color := "#FFFFFF"
	userID := uuid.New()

	// 2. Act
	label, err := labelService.CreateLabel(ctx, name, color, userID)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, label)
	assert.EqualError(t, err, "необходимо указать название метки")

	mockRepo.AssertExpectations(t)
}

func TestCreateLabel_InvalidColor(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	name := "Test Label"
	color := "invalid-color"
	userID := uuid.New()

	// 2. Act
	label, err := labelService.CreateLabel(ctx, name, color, userID)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, label)
	assert.EqualError(t, err, "неверный формат цвета (HEX)")

	mockRepo.AssertExpectations(t)
}

func TestGetLabelByID(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	labelID := uuid.New()
	expectedLabel := &domain.Label{
		ID:     labelID,
		Name:   "Test Label",
		Color:  "#FFFFFF",
		UserID: uuid.New(),
	}

	// Настройка mock-репозитория
	mockRepo.On("GetByID", mock.Anything, labelID).Return(expectedLabel, nil)

	// 2. Act
	label, err := labelService.GetLabelByID(ctx, labelID)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedLabel, label)

	mockRepo.AssertExpectations(t)
}

func TestGetLabelByID_NotFound(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	labelID := uuid.New()

	// Настройка mock-репозитория
	mockRepo.On("GetByID", mock.Anything, labelID).Return(nil, errors.New("label not found"))

	// 2. Act
	label, err := labelService.GetLabelByID(ctx, labelID)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, label)
	assert.EqualError(t, err, "ошибка при получении метки по ID: label not found")

	mockRepo.AssertExpectations(t)
}

func TestUpdateLabel(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	labelID := uuid.New()
	initialLabel := &domain.Label{
		ID:     labelID,
		Name:   "Old Label",
		Color:  "#000000",
		UserID: uuid.New(),
	}
	updatedName := "New Label"
	updatedColor := "#FFFFFF"

	// Настройка mock-репозитория
	mockRepo.On("GetByID", mock.Anything, labelID).Return(initialLabel, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(label *domain.Label) bool {
		return label.ID == labelID && label.Name == updatedName && label.Color == updatedColor
	})).Return(nil)

	// 2. Act
	label, err := labelService.UpdateLabel(ctx, labelID, updatedName, updatedColor)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, labelID, label.ID)
	assert.Equal(t, updatedName, label.Name)
	assert.Equal(t, updatedColor, label.Color)

	mockRepo.AssertExpectations(t)
}

func TestUpdateLabel_NotFound(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	labelID := uuid.New()
	updatedName := "New Label"
	updatedColor := "#FFFFFF"

	// Настройка mock-репозитория
	mockRepo.On("GetByID", mock.Anything, labelID).Return(nil, errors.New("label not found"))

	// 2. Act
	label, err := labelService.UpdateLabel(ctx, labelID, updatedName, updatedColor)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, label)
	assert.EqualError(t, err, "метка не найдена")

	mockRepo.AssertExpectations(t)
}

func TestDeleteLabel(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	labelID := uuid.New()

	// Настройка mock-репозитория
	mockRepo.On("Delete", mock.Anything, labelID).Return(nil)

	// 2. Act
	err := labelService.DeleteLabel(ctx, labelID)

	// 3. Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteLabel_Error(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	labelID := uuid.New()
	expectedError := errors.New("delete error")

	// Настройка mock-репозитория
	mockRepo.On("Delete", mock.Anything, labelID).Return(expectedError)

	// 2. Act
	err := labelService.DeleteLabel(ctx, labelID)

	// 3. Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "ошибка при удалении метки: delete error")

	mockRepo.AssertExpectations(t)
}

func TestGetAllLabelsByUserID(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	expectedLabels := []*domain.Label{
		{ID: uuid.New(), Name: "Label 1", Color: "#FFFFFF", UserID: userID},
		{ID: uuid.New(), Name: "Label 2", Color: "#000000", UserID: userID},
	}

	// Настройка mock-репозитория
	mockRepo.On("GetAllByUserID", mock.Anything, userID).Return(expectedLabels, nil)

	// 2. Act
	labels, err := labelService.GetAllLabelsByUserID(ctx, userID)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedLabels, labels)

	mockRepo.AssertExpectations(t)
}

func TestGetAllLabelsByUserID_Error(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockLabelRepository)
	labelService := NewLabelService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	expectedError := errors.New("get all error")

	// Настройка mock-репозитория
	mockRepo.On("GetAllByUserID", mock.Anything, userID).Return(nil, expectedError)

	// 2. Act
	labels, err := labelService.GetAllLabelsByUserID(ctx, userID)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, labels)
	assert.EqualError(t, err, "ошибка при получении меток пользователя: get all error")

	mockRepo.AssertExpectations(t)
}
