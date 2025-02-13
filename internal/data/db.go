package data

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "shop"
)

// Init database connection
func InitDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
}

// Closing database
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
