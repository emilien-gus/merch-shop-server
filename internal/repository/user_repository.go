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
	query := "SELECT id, username, password, balance FROM users WHERE username = $1"

	err := ur.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Balance)
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

	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err = ur.db.Exec(query, username, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetUserInfo(userID int) (*models.InfoResponse, error) {
	userInfo, err := ur.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	inventory, err := ur.getUserInventory(userID)
	if err != nil {
		return nil, err
	}

	transactions, err := ur.getUserTransactions(userID)
	if err != nil {
		return nil, err
	}

	return &models.InfoResponse{
		CoinsCount:  userInfo.Balance,
		Inventory:   inventory,
		CoinHistory: transactions,
	}, nil
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

func (ur *UserRepository) getUserInventory(userID int) ([]models.Item, error) {
	query := "SELECT item_name, quantity FROM purchases WHERE user_id = $1"
	rows, err := ur.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, err
		}
		inventory = append(inventory, item)
	}
	return inventory, nil
}

func (ur *UserRepository) getUserTransactions(userID int) (models.CoinHistory, error) {
	var history models.CoinHistory

	sentQuery := "SELECT sender_id, amount FROM transactions WHERE receiver_id = $1"
	receivedQuery := "SELECT receiver_id, amount FROM transactions WHERE sender_id = $1"

	receivedRows, err := ur.db.Query(receivedQuery, userID)
	if err != nil {
		return history, err
	}
	defer receivedRows.Close()

	for receivedRows.Next() {
		var transaction models.ReceivedTransaction
		var senderID int
		if err := receivedRows.Scan(&senderID, &transaction.Amount); err != nil {
			return history, err
		}
		transaction.FromUser, _ = ur.getUsernameByID(senderID)
		history.Received = append(history.Received, transaction)
	}

	sentRows, err := ur.db.Query(sentQuery, userID)
	if err != nil {
		return history, err
	}
	defer sentRows.Close()

	for sentRows.Next() {
		var transaction models.SentTransaction
		var receiverID int
		if err := sentRows.Scan(&receiverID, &transaction.Amount); err != nil {
			return history, err
		}
		transaction.ToUser, _ = ur.getUsernameByID(receiverID)
		history.Sent = append(history.Sent, transaction)
	}

	return history, nil
}

func (r *UserRepository) getUsernameByID(userID int) (string, error) {
	var username string
	query := "SELECT username FROM users WHERE id = $1"
	err := r.db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
