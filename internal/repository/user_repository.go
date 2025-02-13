package repository

import (
	"avito-shop/internal/models"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, balance FROM users WHERE username = $1"

	err := ur.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Если пользователь не найден, возвращаем nil
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, balance FROM users WHERE id = $1"

	err := ur.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Если пользователь не найден, возвращаем nil
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) InsertUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// SQL-запрос на вставку
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err = ur.db.Exec(query, username, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}
