package usecase

import (
	"fmt"
	"go.uber.org/zap"
	"mnc-test/model"
	"mnc-test/repository"
)

type TransactionUsecase interface {
	MakePayment(tx *model.Transaction) error
	GetCustomerTransaction(custID int) ([]model.Transaction, error)
}

type transactionUsecase struct {
	txRepo repository.TransactionRepository
	logger *zap.Logger
}

func (t transactionUsecase) MakePayment(tx *model.Transaction) error {
	//TODO implement me

	if tx.MerchantID == 0 {
		return fmt.Errorf("merchant ID cannot be zero")
	} else if tx.Amount == 0 {
		return fmt.Errorf("amount cannot be zero")
	}

	if err := t.txRepo.MakePayment(tx); err != nil {
		return fmt.Errorf(err.Error())
	}

	//log
	if t.logger != nil {
		t.logger.Info("A payment has been made",
			zap.Int("customerID", tx.CustomerID),
			zap.Int("merchantID", tx.MerchantID),
			zap.Int("amount", tx.Amount))
	} else {
		fmt.Println("Logger is not initialized")
	}

	return nil
}

func (t transactionUsecase) GetCustomerTransaction(custID int) ([]model.Transaction, error) {
	//TODO implement me

	if custID == 0 {
		return nil, fmt.Errorf("customer id cannot empty")
	}

	txs, err := t.txRepo.GetCustomerTransactionByID(custID)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return txs, nil
}

func NewTransactionUsecase(txrepo repository.TransactionRepository, logger *zap.Logger) TransactionUsecase {
	return &transactionUsecase{
		txRepo: txrepo,
		logger: logger,
	}
}
