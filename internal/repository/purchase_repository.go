package repository

import (
	"database/sql"
	"errors"
)

type PurchaseRepository struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (pr *PurchaseRepository) InsertBuying(userID int, item string, price int) error {
	tx, err := pr.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Вычитаем монеты из баланса
	result, err := tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2 AND balance >= $1", price, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("insufficient funds")
	}

	result, err = tx.Exec("UPDATE purchases SET quantity = quantity + 1 WHERE user_id = $1 AND item_name = $2", userID, item)
	if err != nil {
		return err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		_, err = tx.Exec("INSERT INTO purchases (user_id, item_name, price) VALUES ($1, $2, $3)", userID, item, price)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
