package usecase

import (
	"fmt"
	"mnc-test/model"
	"mnc-test/repository"
	"mnc-test/util/helper"
)

type TransactionUsecase interface {
	MakePayment(tx *model.Transaction) error
	GetCustomerTransaction(custID int) ([]model.Transaction, error)
}

type transactionUsecase struct {
	txRepo repository.TransactionRepository
}

func (t transactionUsecase) MakePayment(tx *model.Transaction) error {
	//TODO implement me

	if tx.MerchantID == 0 {
		return fmt.Errorf("merchant ID cannot be zero")
	} else if tx.Amount == 0 {
		return fmt.Errorf("amount cannot be zero")
	}

	//generate uuid for transaction id
	tx.TransactionID = helper.GenerateUUID()

	if err := t.txRepo.MakePayment(tx); err != nil {
		return fmt.Errorf(err.Error())
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

func NewTransactionUsecase(txrepo repository.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{
		txRepo: txrepo,
	}
}
