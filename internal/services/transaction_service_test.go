package services

import (
	"avito-shop/internal/models"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockTransactionRepository struct {
	mock.Mock
	users map[string]*models.User // Сохраняем пользователей в мапе
}

// Конструктор мока
func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		users: make(map[string]*models.User), // Инициализируем пустую карту пользователей
	}
}

// Метод для вставки пользователя (например, при регистрации)
func (m *MockTransactionRepository) AddUser(username string, id int, password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	currBalance := 1000
	m.users[username] = &models.User{
		ID:       id,
		Username: username,
		Password: string(hashedPassword),
		Balance:  currBalance,
	}
}

// Метод для получения ID пользователя по имени
func (m *MockTransactionRepository) GetUserIDByUsername(username string) (int, error) {
	// Проверяем, есть ли такой пользователь в моке
	if user, exists := m.users[username]; exists {
		return user.ID, nil // Возвращаем ID, если пользователь найден
	}
	return 0, errors.New("user not found") // Ошибка, если пользователь не найден
}

func (m *MockTransactionRepository) InsertTransaction(senderID int, receiver int, amount int) error {
	// Мокаем создание транзакции
	args := m.Called(senderID, receiver, amount)
	return args.Error(0)
}

func TestSendCoins_Success(t *testing.T) {
	mockRepo := NewMockTransactionRepository()
	service := NewTransactionService(mockRepo)

	// Добавляем пользователей в мок
	mockRepo.AddUser("sender_user", 1, "password")
	mockRepo.AddUser("receiver_user", 2, "password")

	// Мокаем получение ID получателя

	// Мокаем успешное создание транзакции
	mockRepo.On("InsertTransaction", 1, 2, 100).Return(nil).Once()

	// Вызываем метод
	err := service.SendCoins(1, "receiver_user", 100)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Проверяем, что все ожидания были выполнены
	mockRepo.AssertExpectations(t)
}

func TestSendCoins_UserNotFound(t *testing.T) {
	mockRepo := NewMockTransactionRepository()
	service := NewTransactionService(mockRepo)

	// Добавляем отправителя в мок
	mockRepo.AddUser("sender_user", 1, "password")

	// Проверяем, что метод вернет ошибку "user not found"
	err := service.SendCoins(1, "non_existent_user", 100)

	// Проверяем, что ошибка соответствует ожидаемой
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	// Проверяем, что моки были вызваны с правильными аргументами
	mockRepo.AssertExpectations(t)
}

func TestSendCoins_AmountError(t *testing.T) {
	mockRepo := NewMockTransactionRepository()
	service := NewTransactionService(mockRepo)

	// Проверяем, что метод вернет ошибку "user not found"
	err := service.SendCoins(1, "non_existent_user", 0)

	// Проверяем, что ошибка соответствует ожидаемой
	assert.Error(t, err)
	assert.Equal(t, "amount must be greater than zero", err.Error())

	// Проверяем, что моки были вызваны с правильными аргументами
	mockRepo.AssertExpectations(t)
}

func TestSendCoins_LowAmountError(t *testing.T) {
	mockRepo := NewMockTransactionRepository()
	service := NewTransactionService(mockRepo)

	mockRepo.AddUser("sender_user", 1, "password")
	mockRepo.AddUser("receiver_user", 2, "password")

	mockRepo.On("InsertTransaction", 1, 2, 1001).Return(fmt.Errorf("low balance")).Once()

	// Проверяем, что метод вернет ошибку "user not found"
	err := service.SendCoins(1, "receiver_user", 1001)

	// Проверяем, что ошибка соответствует ожидаемой
	assert.Error(t, err)
	assert.Equal(t, "low balance", err.Error())

	// Проверяем, что моки были вызваны с правильными аргументами
	mockRepo.AssertExpectations(t)
}
