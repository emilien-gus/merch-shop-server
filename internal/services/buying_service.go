package services

import "avito-shop/internal/repository"

type BuyingService struct {
	purchaseRepo *repository.PurchaseRepository
}

func NewBuyingService(purchaseRepo *repository.PurchaseRepository) *BuyingService {
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
		return err
	}

	return nil
}
