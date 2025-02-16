package services

import (
	"avito-shop/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	user, _ := args.Get(0).(*models.User) // Если nil, не падаем
	return user, args.Error(1)
}

func (m *MockUserRepository) InsertUser(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0) // Без лишней логики
}

func TestAuthenticateUser_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := NewUserService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &models.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
	}

	// Задаем ожидаемое поведение мока
	mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)

	// Вызываем тестируемую функцию
	token, err := userService.AuthenticateUser("testuser", "password123")

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t) // Проверяем, что мок вызвался корректно
}

func TestAuthenticateUser_NewUserCreated(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := NewUserService(mockRepo)

	// 1. Первый вызов: пользователь не найден
	mockRepo.On("GetUserByUsername", "newuser").Return(nil, nil).Once()

	// 2. Создание пользователя
	mockRepo.On("InsertUser", "newuser", "newpassword").Return(nil).Once()

	// 3. Второй вызов: теперь пользователь найден
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("newpassword"), bcrypt.DefaultCost)
	mockUser := &models.User{
		ID:       2,
		Username: "newuser",
		Password: string(hashedPassword),
	}
	mockRepo.On("GetUserByUsername", "newuser").Return(mockUser, nil).Once()

	// Запускаем тестируемый метод
	token, err := userService.AuthenticateUser("newuser", "newpassword")

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Убеждаемся, что все методы были вызваны в нужном порядке
	mockRepo.AssertExpectations(t)
}
