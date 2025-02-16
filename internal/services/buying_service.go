package services

import (
	"avito-shop/internal/repository"
	"log"
)

type BuyingService struct {
	purchaseRepo repository.PurchaseRepositoryInteface
}

func NewBuyingService(purchaseRepo repository.PurchaseRepositoryInteface) *BuyingService {
	return &BuyingService{purchaseRepo: purchaseRepo}
}

func (s *BuyingService) BuyItem(userId int, item string) error {
	var price int
	price, err := GetItem(item)
	if err != nil {
		return err
	}

	err = s.purchaseRepo.InsertBuying(userId, item, price)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	return nil
}
