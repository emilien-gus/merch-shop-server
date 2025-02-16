package services

import (
	"avito-shop/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository представляет собой мок репозитория пользователей
type MockInfoRepository struct {
	mock.Mock
}

// GetUserInfo мокаем получение информации о пользователе
func (m *MockUserRepository) GetUserInfo(userID int) (*models.InfoResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.InfoResponse), args.Error(1)
}

// Тест для метода GetInfo (успешный случай)
func TestGetInfo_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewInfoService(mockRepo)

	// Мокаем успешное возвращение информации о пользователе
	mockRepo.On("GetUserInfo", 1).Return(&models.InfoResponse{
		CoinsCount:  1000,
		Inventory:   []models.Item{{Type: "t-shirt", Quantity: 1}},
		CoinHistory: models.CoinHistory{},
	}, nil).Once()

	// Вызываем метод
	info, err := service.GetInfo(1)

	// Проверяем, что ошибки нет и информация возвращена корректно
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, 1000, info.CoinsCount)
	mockRepo.AssertExpectations(t)
}

// Тест для метода GetInfo (ошибка при получении данных о пользователе)
func TestGetInfo_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewInfoService(mockRepo)

	// Мокаем ошибку при получении информации о пользователе
	mockRepo.On("GetUserInfo", 1).Return((*models.InfoResponse)(nil), errors.New("user not found")).Once()

	// Вызываем метод
	info, err := service.GetInfo(1)

	// Проверяем, что ошибка была возвращена
	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
