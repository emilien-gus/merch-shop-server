package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

const secret = "w8jZ6FcR2vLp7XQ1mKtY9uHnJ3bGq4As5Df0eVxIyNrO"

func setupTestDB() *sql.DB {
	// Создаем подключение к тестовой базе данных, например, PostgreSQL или SQLite
	connStr := "host=postgres port=5432 user=postgres password=password dbname=shop sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
		return nil
	}

	// Выполняем миграции или создаем таблицы для теста
	return db
}

func authenticateUser(baseURL, username, password string) string {
	authRequest := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(authRequest)
	resp, _ := http.Post(baseURL+"/auth", "application/json", bytes.NewBuffer(body))

	var authResponse map[string]string
	json.NewDecoder(resp.Body).Decode(&authResponse)
	return authResponse["token"]
}

func deleteUserByID(db *sql.DB, userID int) error {
	// Подготовленный запрос для защиты от SQL-инъекций
	query := "DELETE FROM users WHERE id = $1"
	_, err := db.Exec(query, userID)
	return err
}
