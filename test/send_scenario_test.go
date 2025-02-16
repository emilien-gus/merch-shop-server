package test

import (
	"avito-shop/internal/handlers"
	"avito-shop/internal/middleware"
	"avito-shop/internal/models"
	"avito-shop/internal/repository"
	"avito-shop/internal/services"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupRouterSend(transactionHandler *handlers.TransactionHandler, authHandler *handlers.AuthHandler, middleware gin.HandlerFunc) *gin.Engine {

	router := gin.Default()
	router.Group("/api")
	router.POST("/auth", authHandler.Login)
	router.Use(middleware)
	router.POST("/sendCoins", transactionHandler.SendCoin)
	return router
}

func TestSendCoins(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	middleware.SetSecretKey(secret)

	repo := repository.NewTransactionRepository(db)
	authRepo := repository.NewUserRepository(db)

	service := services.NewTransactionService(repo)
	authServ := services.NewUserService(authRepo)

	handler := handlers.NewTransactionHandler(service)
	auth := handlers.NewAuthHandler(authServ)

	sender := models.User{ID: 2, Username: "send_user", Balance: 1000}
	receiver := models.User{ID: 3, Username: "receiver_user", Balance: 1000}

	router := setupRouterSend(handler, auth, middleware.JWTMiddleware())
	server := httptest.NewServer(router)
	defer server.Close()

	token := authenticateUser(server.URL, sender.Username, "password")
	authenticateUser(server.URL, receiver.Username, "password")

	requestBody := struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}{
		ToUser: receiver.Username,
		Amount: 500,
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", server.URL+"/sendCoins", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedSender models.User
	var updatedReceiver models.User
	err = db.QueryRow("SELECT username, balance FROM users WHERE id = $1", sender.ID).Scan(&updatedSender.Username, &updatedSender.Balance)
	if err != nil {
		t.Fatalf("Failed to query updated user: %v", err)
	}
	err = db.QueryRow("SELECT username, balance FROM users WHERE id = $1", receiver.ID).Scan(&updatedReceiver.Username, &updatedReceiver.Balance)
	if err != nil {
		t.Fatalf("Failed to query updated user: %v", err)
	}

	err = deleteUserByID(db, updatedSender.ID)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = deleteUserByID(db, updatedReceiver.ID)
	if err != nil {
		log.Fatal(err.Error())
	}
	senderExpectedBalance := sender.Balance - 500
	assert.Equal(t, senderExpectedBalance, updatedSender.Balance)
	assert.Equal(t, sender.Username, updatedSender.Username)

	receiverExpectedBalance := receiver.Balance + 500
	assert.Equal(t, receiverExpectedBalance, updatedReceiver.Balance)
	assert.Equal(t, receiver.Username, updatedReceiver.Username)

}
