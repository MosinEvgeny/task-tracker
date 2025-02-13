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

// MockUserRepository - это mock для UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	user, ok := args.Get(0).(*domain.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	user, ok := args.Get(0).(*domain.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	username := "testuser"
	email := "test@example.com"
	password := "password"

	// Ожидаемый пользователь
	expectedUser := &domain.User{
		Username: username,
		Email:    email,
	}

	// Настройка mock-репозитория
	mockRepo.On("GetByEmail", ctx, email).Return(nil, errors.New("user not found")) // Имитируем, что пользователя с таким email не существует
	mockRepo.On("Create", ctx, mock.MatchedBy(func(user *domain.User) bool {
		expectedUser.ID = user.ID // Сохраняем ID, чтобы потом сравнить
		return user.Username == username && user.Email == email
	})).Return(nil)

	// 2. Act
	user, err := userService.CreateUser(ctx, username, email, password)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NotEmpty(t, user.Password) // Проверяем, что пароль был захэширован

	mockRepo.AssertExpectations(t) // Проверяем, что все ожидаемые вызовы mock-методов были выполнены
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	username := "testuser"
	email := "invalid-email"
	password := "password"

	// 2. Act
	user, err := userService.CreateUser(ctx, username, email, password)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "неверный формат email")

	mockRepo.AssertExpectations(t) // Проверяем, что mock-методы не вызывались
}

func TestCreateUser_ExistingEmail(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	username := "testuser"
	email := "test@example.com"
	password := "password"

	// Настройка mock-репозитория
	mockRepo.On("GetByEmail", ctx, email).Return(&domain.User{ID: uuid.New(), Username: "existingUser", Email: email}, nil) // Имитируем, что пользователь с таким email уже существует

	// 2. Act
	user, err := userService.CreateUser(ctx, username, email, password)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "пользователь с таким email уже существует")

	mockRepo.AssertExpectations(t) // Проверяем, что вызовы mock-методов соответствуют ожиданиям
}

func TestGetUserByID(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	expectedUser := &domain.User{ID: userID, Username: "testuser", Email: "test@example.com"}

	// Настройка mock-репозитория
	mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

	// 2. Act
	user, err := userService.GetUserByID(ctx, userID)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()

	// Настройка mock-репозитория
	mockRepo.On("GetByID", ctx, userID).Return(nil, errors.New("user not found"))

	// 2. Act
	user, err := userService.GetUserByID(ctx, userID)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "ошибка при получении пользователя по ID: user not found")

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	initialUser := &domain.User{ID: userID, Username: "olduser", Email: "old@example.com"}
	updatedUsername := "newuser"
	updatedEmail := "new@example.com"

	// Настройка mock-репозитория
	mockRepo.On("GetByID", ctx, userID).Return(initialUser, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(user *domain.User) bool {
		return user.ID == userID && user.Username == updatedUsername && user.Email == updatedEmail
	})).Return(nil)

	// 2. Act
	user, err := userService.UpdateUser(ctx, userID, updatedUsername, updatedEmail)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, updatedUsername, user.Username)
	assert.Equal(t, updatedEmail, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_NotFound(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	updatedUsername := "newuser"
	updatedEmail := "new@example.com"

	// Настройка mock-репозитория
	mockRepo.On("GetByID", ctx, userID).Return(nil, errors.New("user not found"))

	// 2. Act
	user, err := userService.UpdateUser(ctx, userID, updatedUsername, updatedEmail)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "пользователь не найден")

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()

	// Настройка mock-репозитория
	mockRepo.On("Delete", ctx, userID).Return(nil)

	// 2. Act
	err := userService.DeleteUser(ctx, userID)

	// 3. Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Error(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	expectedError := errors.New("delete error")

	// Настройка mock-репозитория
	mockRepo.On("Delete", ctx, userID).Return(expectedError)

	// 2. Act
	err := userService.DeleteUser(ctx, userID)

	// 3. Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "ошибка при удалении пользователя: delete error")

	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	email := "test@example.com"
	expectedUser := &domain.User{ID: uuid.New(), Username: "testuser", Email: email}

	// Настройка mock-репозитория
	mockRepo.On("GetByEmail", ctx, email).Return(expectedUser, nil)

	// 2. Act
	user, err := userService.GetUserByEmail(ctx, email)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	// 1. Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	ctx := context.Background()

	email := "test@example.com"

	// Настройка mock-репозитория
	mockRepo.On("GetByEmail", ctx, email).Return(nil, errors.New("user not found"))

	// 2. Act
	user, err := userService.GetUserByEmail(ctx, email)

	// 3. Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "ошибка при получении пользователя по email: user not found")

	mockRepo.AssertExpectations(t)
}
