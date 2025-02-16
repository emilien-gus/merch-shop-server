package services

import (
	"avito-shop/internal/repository"
	"errors"
)

type TransactionService struct {
	transactionRepo repository.TransactionRepositoryInterface
}

func NewTransactionService(transactionRepo repository.TransactionRepositoryInterface) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}

func (ts *TransactionService) SendCoins(senderID int, receiverUsername string, amount int) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	receiverID, err := ts.transactionRepo.GetUserIDByUsername(receiverUsername)
	if err != nil {
		return err
	}

	err = ts.transactionRepo.InsertTransaction(senderID, receiverID, amount)
	if err != nil {
		return err
	}

	return nil
}
