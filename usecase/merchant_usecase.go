package usecase

import (
	"fmt"
	"mnc-test/model"
	"mnc-test/repository"
)

type MerchantUsecase interface {
	CreateNewMerchant(merchant model.Merchant) error
	FindAllMerchant() ([]model.Merchant, error)
}

type merchantUsecase struct {
	merchantRepo repository.MerchantRepository
}

func (m *merchantUsecase) CreateNewMerchant(merchant model.Merchant) error {
	//TODO implement me

	if merchant.Name == "" {
		return fmt.Errorf("Name is required")
	}

	if err := m.merchantRepo.CreateNewMerchant(merchant); err != nil {
		return fmt.Errorf("Failed to create merchant %v", err.Error())
	}

	return nil
}

func (m *merchantUsecase) FindAllMerchant() ([]model.Merchant, error) {
	//TODO implement me

	merchant, err := m.merchantRepo.FindAllMerchant()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return merchant, nil
}

func NewMerchantUsecase(muc repository.MerchantRepository) MerchantUsecase {
	return &merchantUsecase{
		merchantRepo: muc,
	}
}
