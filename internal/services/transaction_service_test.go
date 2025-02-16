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
	users map[string]*models.User
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		users: make(map[string]*models.User),
	}
}

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

func (m *MockTransactionRepository) GetUserIDByUsername(username string) (int, error) {
	if user, exists := m.users[username]; exists {
		return user.ID, nil
	}
	return 0, errors.New("user not found")
}

func (m *MockTransactionRepository) InsertTransaction(senderID int, receiver int, amount int) error {
	args := m.Called(senderID, receiver, amount)
	return args.Error(0)
}

func TestSendCoins_Success(t *testing.T) {
	mockRepo := NewMockTransactionRepository()
	service := NewTransactionService(mockRepo)

	mockRepo.AddUser("sender_user", 1, "password")
	mockRepo.AddUser("receiver_user", 2, "password")

	mockRepo.On("InsertTransaction", 1, 2, 100).Return(nil).Once()

	err := service.SendCoins(1, "receiver_user", 100)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSendCoins_UserNotFound(t *testing.T) {
	mockRepo := NewMockTransactionRepository()
	service := NewTransactionService(mockRepo)

	mockRepo.AddUser("sender_user", 1, "password")

	err := service.SendCoins(1, "non_existent_user", 100)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestSendCoins_AmountError(t *testing.T) {
	mockRepo := NewMockTransactionRepository()
	service := NewTransactionService(mockRepo)

	err := service.SendCoins(1, "non_existent_user", 0)

	assert.Error(t, err)
	assert.Equal(t, "amount must be greater than zero", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestSendCoins_LowAmountError(t *testing.T) {
	mockRepo := NewMockTransactionRepository()
	service := NewTransactionService(mockRepo)

	mockRepo.AddUser("sender_user", 1, "password")
	mockRepo.AddUser("receiver_user", 2, "password")

	mockRepo.On("InsertTransaction", 1, 2, 1001).Return(fmt.Errorf("low balance")).Once()

	err := service.SendCoins(1, "receiver_user", 1001)

	assert.Error(t, err)
	assert.Equal(t, "low balance", err.Error())

	mockRepo.AssertExpectations(t)
}
