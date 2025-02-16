package handlers

import (
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	infoService *services.InfoService
}

func NewInfoHandler(infoService *services.InfoService) *InfoHandler {
	return &InfoHandler{infoService: infoService}
}

func (ih *InfoHandler) UserInfo(c *gin.Context) {
	userId, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	info, err := ih.infoService.GetInfo(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
