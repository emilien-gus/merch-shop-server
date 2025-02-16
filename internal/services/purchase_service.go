package services

import (
	"avito-shop/internal/repository"
	"log"
)

type PurchaseService struct {
	purchaseRepo repository.PurchaseRepositoryInteface
}

func NewPurchaseService(purchaseRepo repository.PurchaseRepositoryInteface) *PurchaseService {
	return &PurchaseService{purchaseRepo: purchaseRepo}
}

func (s *PurchaseService) BuyItem(userId int, item string) error {
	var price int
	price, err := GetItem(item)
	if err != nil {
		return err
	}

	err = s.purchaseRepo.InsertPurchase(userId, item, price)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	return nil
}
