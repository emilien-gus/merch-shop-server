package repository

import (
	"database/sql"
	"errors"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (tr *TransactionRepository) InsertTransaction(senderID int, receiver string, amount int) error {
	tx, err := tr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	receiverID, err := tr.GetUserIDByUsername(receiver)
	if err != nil {
		return err
	}

	res, err := tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, senderID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("insufficient funds")
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", amount, receiverID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions (sender_id, receiver_id, amount) VALUES ($1, $2, $3)", senderID, receiverID, amount)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (tr *TransactionRepository) GetUserIDByUsername(username string) (int, error) {
	var userID int
	err := tr.db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return userID, nil
}
