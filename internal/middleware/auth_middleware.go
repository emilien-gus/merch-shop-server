package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var SecretKey []byte

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Загружаем секретный ключ из переменной окружения
		if string(SecretKey) == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: JWT secret is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем алгоритм подписи токена
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return SecretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Проверяем claims и устанавливаем user_id в контекст
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userID, exists := claims["user_id"]; exists {
				c.Set("user_id", userID)
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: missing user_id"})
			c.Abort()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
	}
}

func InitSecretKey() {
	SecretKey = []byte(os.Getenv("JWT_SECRET"))
}
