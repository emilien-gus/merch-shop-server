package services

import (
	"avito-shop/internal/models"
	"avito-shop/internal/repository"
)

type InfoService struct {
	userRepo repository.UserInfoRepositoryInterfase
}

func NewInfoService(userRepo repository.UserInfoRepositoryInterfase) *InfoService {
	return &InfoService{userRepo: userRepo}
}

func (is *InfoService) GetInfo(userID int) (*models.InfoResponse, error) {
	info, err := is.userRepo.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
