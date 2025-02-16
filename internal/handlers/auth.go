package handlers

import (
	"avito-shop/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Проверяем логин и пароль
	token, err := ah.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, errors.New("wrong password")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Отправляем токен клиенту
	c.JSON(http.StatusOK, gin.H{"token": token})
}
