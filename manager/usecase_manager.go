package manager

import (
	"go.uber.org/zap"
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
	rm  RepoManager
	log *zap.Logger
}

func (u usecaseManager) UserCredUsecase() usecase.UserCredentialUsecase {
	return usecase.NewUserCredentialUsecase(u.rm.UserCredRepo(), u.log)
}

func (u usecaseManager) MerchantUsecase() usecase.MerchantUsecase {
	//TODO implement me
	return usecase.NewMerchantUsecase(u.rm.MerchantRepo())
}

func (u usecaseManager) TransactionUsecase() usecase.TransactionUsecase {
	//TODO implement me
	return usecase.NewTransactionUsecase(u.rm.TransactionRepo(), u.log)
}

func (u usecaseManager) TransferUsecase() usecase.TransferUsecase {
	//TODO implement me
	return usecase.NewTransferUsecase(u.rm.TransferRepo(), u.log)
}

func NewUsecaseManager(rm RepoManager, log *zap.Logger) UsecaseManager {
	return &usecaseManager{
		rm:  rm,
		log: log,
	}
}
