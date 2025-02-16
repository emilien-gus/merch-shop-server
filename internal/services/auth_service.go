package services

import (
	// Импортируем секретный ключ
	"avito-shop/internal/models"
	"avito-shop/internal/repository"
	"os"

	"errors"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepositoryInterface
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

// AuthenticateUser проверяет учетные данные и возвращает JWT
func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if user == nil {
		err = s.userRepo.InsertUser(username, password)
		if err != nil {
			return "", err
		}
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
		return "", err
	}

	return token, nil
}

type CustomClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// generateJWT создает JWT-токен для пользователя
func (s *UserService) generateJWT(user *models.User) (string, error) {
	claims := CustomClaims{
		UserID: user.ID,
	}

	claims.UserID = user.ID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
