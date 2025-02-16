package services

import (
	"avito-shop/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockInfoRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserInfo(userID int) (*models.InfoResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.InfoResponse), args.Error(1)
}

func TestGetInfo_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewInfoService(mockRepo)

	mockRepo.On("GetUserInfo", 1).Return(&models.InfoResponse{
		CoinsCount:  1000,
		Inventory:   []models.Item{{Type: "t-shirt", Quantity: 1}},
		CoinHistory: models.CoinHistory{},
	}, nil).Once()

	info, err := service.GetInfo(1)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, 1000, info.CoinsCount)
	mockRepo.AssertExpectations(t)
}

func TestGetInfo_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewInfoService(mockRepo)

	mockRepo.On("GetUserInfo", 1).Return((*models.InfoResponse)(nil), errors.New("user not found")).Once()

	info, err := service.GetInfo(1)

	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
