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
	transactRepo := repository.NewTransactionRepository(db)

	userService := services.NewUserService(userRep)
	purchSer := services.NewPurchaseService(purchRep)
	transactSer := services.NewTransactionService(transactRepo)
	infoSer := services.NewInfoService(userRep)

	userHandler := NewAuthHandler(userService)
	purchHandler := NewBuyingHandler(purchSer)
	transactHandler := NewTransactionHandler(transactSer)
	infoHandler := NewInfoHandler(infoSer)

	api := r.Group("/api")

	api.POST("/auth", userHandler.Login)

	api.Use(middleware.JWTMiddleware())
	api.POST("/buy/:item", purchHandler.Buy)
	api.POST("/sendCoin", transactHandler.SendCoin)
	api.GET("/info", infoHandler.UserInfo)
}
