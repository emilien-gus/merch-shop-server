package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPurchaseRepository struct {
	mock.Mock
}

func NewMockPurchaseRepository() *MockPurchaseRepository {
	return &MockPurchaseRepository{}
}

func (m *MockPurchaseRepository) InsertPurchase(userID int, item string, price int) error {
	args := m.Called(userID, item, price)
	return args.Error(0)
}

func TestInsertBuying_Success(t *testing.T) {
	mockRepo := NewMockPurchaseRepository()
	service := NewPurchaseService(mockRepo)

	mockRepo.On("InsertBuying", 1, "t-shirt", 80).Return(nil)

	err := service.BuyItem(1, "t-shirt")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestInsertBuying_NoItemError(t *testing.T) {
	mockRepo := NewMockPurchaseRepository()
	service := NewPurchaseService(mockRepo)

	err := service.BuyItem(1, "ball")

	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
	mockRepo.AssertExpectations(t)
}
