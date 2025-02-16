package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPurchaseRepository представляет собой мок репозитория для покупок
type MockPurchaseRepository struct {
	mock.Mock
}

// NewMockPurchaseRepository создает новый экземпляр мок репозитория для покупок
func NewMockPurchaseRepository() *MockPurchaseRepository {
	return &MockPurchaseRepository{}
}

// InsertBuying мокаем вставку покупки
func (m *MockPurchaseRepository) InsertBuying(userID int, item string, price int) error {
	args := m.Called(userID, item, price)
	return args.Error(0) // Возвращаем ошибку, если она была настроена
}

// Тестирование метода покупки
func TestInsertBuying_Success(t *testing.T) {
	mockRepo := NewMockPurchaseRepository()
	service := NewBuyingService(mockRepo)

	mockRepo.On("InsertBuying", 1, "t-shirt", 80).Return(nil)

	err := service.BuyItem(1, "t-shirt")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestInsertBuying_NoItemError(t *testing.T) {
	mockRepo := NewMockPurchaseRepository()
	service := NewBuyingService(mockRepo)

	err := service.BuyItem(1, "ball")

	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
	mockRepo.AssertExpectations(t)
}
