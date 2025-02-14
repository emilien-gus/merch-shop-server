package handlers

import (
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (th *TransactionHandler) SendCoin(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err = th.transactionService.SendCoins(userID, req.ToUser, req.Amount)
	if err != nil {
		errorString := err.Error()
		if errorString == "insufficient funds" || errorString == "user not found" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "coins send successful"})
}
