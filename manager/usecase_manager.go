package manager

import (
	"github.com/redis/go-redis/v9"
	"mnc-test/usecase"
)

type UsecaseManager interface {
	//all usecase object goes here
	UserCredUsecase() usecase.UserCredentialUsecase
	MerchantUsecase() usecase.MerchantUsecase
	TransactionUsecase() usecase.TransactionUsecase
	TransferUsecase() usecase.TransferUsecase
}

type usecaseManager struct {
	rm     RepoManager
	redisC *redis.Client
}

func (u usecaseManager) UserCredUsecase() usecase.UserCredentialUsecase {
	return usecase.NewUserCredentialUsecase(u.rm.UserCredRepo(), u.redisC)
}

func (u usecaseManager) MerchantUsecase() usecase.MerchantUsecase {
	//TODO implement me
	return usecase.NewMerchantUsecase(u.rm.MerchantRepo())
}

func (u usecaseManager) TransactionUsecase() usecase.TransactionUsecase {
	//TODO implement me
	return usecase.NewTransactionUsecase(u.rm.TransactionRepo())
}

func (u usecaseManager) TransferUsecase() usecase.TransferUsecase {
	//TODO implement me
	return usecase.NewTransferUsecase(u.rm.TransferRepo())
}

func NewUsecaseManager(rm RepoManager, redisC *redis.Client) UsecaseManager {
	return &usecaseManager{
		rm:     rm,
		redisC: redisC,
	}
}
