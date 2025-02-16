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

	mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)

	token, err := userService.AuthenticateUser("testuser", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser_NewUserCreated(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := NewUserService(mockRepo)

	mockRepo.On("GetUserByUsername", "newuser").Return(nil, nil).Once()

	mockRepo.On("InsertUser", "newuser", "newpassword").Return(nil).Once()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("newpassword"), bcrypt.DefaultCost)
	mockUser := &models.User{
		ID:       2,
		Username: "newuser",
		Password: string(hashedPassword),
	}
	mockRepo.On("GetUserByUsername", "newuser").Return(mockUser, nil).Once()

	token, err := userService.AuthenticateUser("newuser", "newpassword")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
}
