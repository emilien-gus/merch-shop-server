package handlers

import (
	"avito-shop/internal/middleware"
	"avito-shop/internal/repository"
	"avito-shop/internal/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB, r *gin.Engine) {
	middleware.InitSecretKey()

	userRep := repository.NewUserRepository(db)
	userService := services.NewUserService(userRep)
	userHandler := NewAuthHandler(userService)
	api := r.Group("/api")

	api.Use(middleware.JWTMiddleware())
	api.POST("/auth", userHandler.Login)
	api.GET("/info")
	api.POST("/sendCoin")
	api.POST("/buy/:item")
}
