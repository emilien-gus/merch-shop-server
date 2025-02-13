package services

import (
	"avito-shop/internal/middleware" // Импортируем секретный ключ
	"avito-shop/internal/models"
	"avito-shop/internal/repository"

	"errors"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// AuthenticateUser проверяет учетные данные и возвращает JWT
func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if user == nil {
		s.userRepo.InsertUser(username, password)
		user, err = s.userRepo.GetUserByUsername(username)
		if err != nil {
			return "", err
		}
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong password")
	}

	// Генерируем JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return "", errors.New("Token generation error")
	}

	return token, nil
}

// generateJWT создает JWT-токен для пользователя
func (s *UserService) generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(middleware.SecretKey)) // Используем секрет из middleware
}
