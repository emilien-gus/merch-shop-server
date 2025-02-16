package handlers

import (
	"avito-shop/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	buyingService *services.PurchaseService
}

func NewBuyingHandler(buyingService *services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{buyingService: buyingService}
}

func (bh *PurchaseHandler) Buy(c *gin.Context) {
	item := c.Param("item")

	userID, err := GetUserID(c)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bh.buyingService.BuyItem(userID, item)
	if err != nil {
		log.Printf("Error: %v", err)
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
