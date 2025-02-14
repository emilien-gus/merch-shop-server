package services

import (
	"avito-shop/internal/repository"
	"errors"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}

func (ts *TransactionService) SendCoins(senderID int, receiverUsername string, amount int) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	err := ts.transactionRepo.InsertTransaction(senderID, receiverUsername, amount)
	if err != nil {
		return err
	}

	return nil
}
