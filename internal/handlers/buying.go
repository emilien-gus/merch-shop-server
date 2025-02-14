package handlers

import (
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BuyingHandler struct {
	buyingService *services.BuyingService
}

func NewBuyingHandler(buyingService *services.BuyingService) *BuyingHandler {
	return &BuyingHandler{buyingService: buyingService}
}

func (bh *BuyingHandler) Buy(c *gin.Context) {
	item := c.Param("item")

	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bh.buyingService.BuyItem(userID, item)
	if err != nil {
		errorString := err.Error()
		if errorString == "item not found" || errorString == "insufficient funds" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "buying item successful"})

}
