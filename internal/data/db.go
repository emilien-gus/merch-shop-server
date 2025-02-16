package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

var DB *sql.DB

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "shop"
)

// Init database connection
func InitDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("Database connection error: %v", err)
	}

	err = DB.Ping() // Пингуем базу данных, чтобы убедиться в подключении
	if err != nil {
		return fmt.Errorf("Failed to connect database: %v", err)
	}
	return nil
}

// Closing database
func CloseDB() error {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("DB is nil at closing")
}
