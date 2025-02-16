package test

import (
	"avito-shop/internal/handlers"
	"avito-shop/internal/middleware"
	"avito-shop/internal/models"
	"avito-shop/internal/repository"
	"avito-shop/internal/services"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupRouterBuy(purchaseHandler *handlers.PurchaseHandler, authHandler *handlers.AuthHandler, middleware gin.HandlerFunc) *gin.Engine {

	router := gin.Default()
	router.Group("/api")
	router.POST("/auth", authHandler.Login)
	router.Use(middleware)
	router.POST("/buy/:item", purchaseHandler.Buy)
	return router
}

func TestPurchaseMerch(t *testing.T) {
	db := setupTestDB()
	defer db.Close()
	middleware.SetSecretKey(secret)

	repo := repository.NewPurchaseRepository(db)
	authRepo := repository.NewUserRepository(db)

	service := services.NewPurchaseService(repo)
	authServ := services.NewUserService(authRepo)

	handler := handlers.NewBuyingHandler(service)
	auth := handlers.NewAuthHandler(authServ)

	user := models.User{ID: 1, Username: "testuser", Balance: 1000}
	item := "t-shirt"

	router := setupRouterBuy(handler, auth, middleware.JWTMiddleware())
	server := httptest.NewServer(router)
	defer server.Close()

	token := authenticateUser(server.URL, "testuser", "password")
	log.Println("JWT token is: " + token)

	req, _ := http.NewRequest("POST", server.URL+"/buy/"+item, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 7. Проверка изменений в базе данных
	var updatedUser models.User
	err = db.QueryRow("SELECT username, balance FROM users WHERE id = $1", user.ID).Scan(&updatedUser.Username, &updatedUser.Balance)
	if err != nil {
		t.Fatalf("Failed to query updated user: %v", err)
	}

	err = deleteUserByID(db, user.ID)
	if err != nil {
		log.Fatal(err.Error())
	}

	expectedBalance := user.Balance - 80
	assert.Equal(t, expectedBalance, updatedUser.Balance)
	assert.Equal(t, user.Username, updatedUser.Username)
}
