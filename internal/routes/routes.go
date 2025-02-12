package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/auth")
	api.GET("/info")
	api.POST("/sendCoin")
	api.POST("/buy/:item")
}
