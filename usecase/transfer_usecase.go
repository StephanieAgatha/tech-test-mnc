package usecase

import (
	"fmt"
	"go.uber.org/zap"
	"mnc-test/model"
	"mnc-test/repository"
	"mnc-test/util/helper"
)

type TransferUsecase interface {
	MakeTransferAccNumbToAccNumb(transactionID string, senderAccountNumber string, receiverAccountNumber string, amount int) error
	GetIncomingMoney(customerId int) ([]model.TransferHistoryIncome, error)
	GetOutcomeMoney(customerId int) ([]model.TransferHistoryOutcome, error)
}

type transferUsecase struct {
	tfRepo repository.TransferRepository
	log    *zap.Logger
}

func (t *transferUsecase) MakeTransferAccNumbToAccNumb(transactionID string, senderAccountNumber string, receiverAccountNumber string, amount int) error {
	//TODO implement me

	//validasi dsini
	if senderAccountNumber == "" {
		return fmt.Errorf("account number cannot be empty")
	} else if receiverAccountNumber == "" {
		return fmt.Errorf("account number cannot be empty")
	} else if amount <= 0 {
		return fmt.Errorf("amount must greater than zero")
	}

	//generate uuid in here
	transactionID = helper.GenerateUUID()

	if err := t.tfRepo.MakeTransferAccNumbToAccNumb(transactionID, senderAccountNumber, receiverAccountNumber, amount); err != nil {
		return fmt.Errorf(err.Error())
	}

	//log
	if t.log != nil {
		t.log.Info("Request transfer money has been initiated",
			zap.String("senderAcountNumber", senderAccountNumber),
			zap.String("receiverAccountNumber", receiverAccountNumber),
			zap.Int("amount", amount))
	} else {
		fmt.Println("Logger is not initialized")
	}

	return nil
}

func (t *transferUsecase) GetIncomingMoney(customerId int) ([]model.TransferHistoryIncome, error) {
	//TODO implement me

	if customerId == 0 {
		return nil, fmt.Errorf("ID is required")
	}

	incomings, err := t.tfRepo.GetIncomingMoney(customerId)
	if err != nil {
		return nil, err
	}

	return incomings, nil
}

func (t *transferUsecase) GetOutcomeMoney(customerId int) ([]model.TransferHistoryOutcome, error) {
	//TODO implement me
	if customerId == 0 {
		return nil, fmt.Errorf("ID is required")
	}

	outcomings, err := t.tfRepo.GetOutcomeMoney(customerId)
	if err != nil {
		return nil, err
	}

	return outcomings, nil
}

func NewTransferUsecase(tfrepo repository.TransferRepository, log *zap.Logger) TransferUsecase {
	return &transferUsecase{
		tfRepo: tfrepo,
		log:    log,
	}
}

//func (t *transferUsecase) MakeTransferPhoneNumbToPhoneNumb(transactionID string, senderPhoneNumber string, receiverPhoneNumber string, amount int) error {
//	//TODO implement me
//	//validasi dsini
//	if senderPhoneNumber == "" {
//		return fmt.Errorf("account number cannot be empty")
//	} else if receiverPhoneNumber == "" {
//		return fmt.Errorf("account number cannot be empty")
//	} else if amount <= 0 {
//		return fmt.Errorf("amount must greater than zero")
//	}
//
//	if err := t.tfRepo.MakeTransferPhoneNumbToPhoneNumb(transactionID, senderPhoneNumber, receiverPhoneNumber, amount); err != nil {
//		return fmt.Errorf(err.Error())
//	}
//
//	//log
//	if t.log != nil {
//		t.log.Info("Request transfer money has been initiated",
//			zap.String("senderPhoneNumber", senderPhoneNumber),
//			zap.String("receiverPhoneNumber", receiverPhoneNumber),
//			zap.Int("amount", amount))
//	} else {
//		fmt.Println("Logger is not initialized")
//	}
//
//	return nil
//}
