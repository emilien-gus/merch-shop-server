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
	purchRep := repository.NewPurchaseRepository(db)

	userService := services.NewUserService(userRep)
	purchSer := services.NewBuyingService(purchRep)

	userHandler := NewAuthHandler(userService)
	purchHandler := NewBuyingHandler(purchSer)

	api := r.Group("/api")

	api.Use(middleware.JWTMiddleware())
	api.POST("/auth", userHandler.Login)
	api.POST("/buy/:item", purchHandler.Buy)
	api.POST("/sendCoin")
	api.GET("/info")
}
