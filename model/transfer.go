package model

import "time"

type TransferRequest struct {
	TransactionID         string `json:"transaction_id"`
	SenderAccountNumber   string `json:"sender_account_number"`
	ReceiverAccountNumber string `json:"receiver_account_number"`
	Amount                int    `json:"amount"`
}

type TransferHistoryIncome struct {
	ID                    int       `json:"id"`
	TransactionID         string    `json:"transaction_id,omitempty"`
	SenderAccountNumber   string    `json:"sender_account_number"`
	ReceiverAccountNumber string    `json:"receiver_account_number"`
	CustomerID            int       `json:"customer_id"`
	Amount                int       `json:"amount"`
	SenderBankName        string    `json:"sender_bank_name"`
	SenderBankId          int       `json:"sender_bank_id"`
	ReceiverBankName      string    `json:"receiver_bank_name"`
	ReceiverBankId        int       `json:"receiver_bank_id"`
	TransferTimeStamp     time.Time `json:"transfer_timestamp"`
}

type TransferHistoryIncomeResponse struct {
	ID                    int    `json:"id"`
	TransactionID         string `json:"transaction_id,omitempty"`
	SenderAccountNumber   string `json:"sender_account_number"`
	ReceiverAccountNumber string `json:"receiver_account_number"`
	Amount                int    `json:"amount"`
	TransferTime          string `json:"transfer_time"`
	SenderBankName        string `json:"sender_bank_name"`
	SenderBankId          int    `json:"sender_bank_id"`
	ReceiverBankName      string `json:"receiver_bank_name"`
	ReceiverBankId        int    `json:"receiver_bank_id"`
}

type TransferHistoryOutcome struct {
	ID                    int       `json:"id"`
	TransactionID         string    `json:"transaction_id,omitempty"`
	SenderAccountNumber   string    `json:"sender_account_number"`
	ReceiverAccountNumber string    `json:"receiver_account_number"`
	CustomerID            int       `json:"customer_id"`
	Amount                int       `json:"amount"`
	SenderBankName        string    `json:"sender_bank_name"`
	SenderBankId          int       `json:"sender_bank_id"`
	ReceiverBankName      string    `json:"receiver_bank_name"`
	ReceiverBankId        int       `json:"receiver_bank_id"`
	TransferTimeStamp     time.Time `json:"transfer_timestamp"`
}

type TransferHistoryOutcomeResponse struct {
	ID                    int    `json:"id"`
	TransactionID         string `json:"transaction_id,omitempty"`
	SenderAccountNumber   string `json:"sender_account_number"`
	ReceiverAccountNumber string `json:"receiver_account_number"`
	Amount                int    `json:"amount"`
	TransferTime          string `json:"transfer_time"`
	SenderBankName        string `json:"sender_bank_name"`
	SenderBankId          int    `json:"sender_bank_id"`
	ReceiverBankName      string `json:"receiver_bank_name"`
	ReceiverBankId        int    `json:"receiver_bank_id"`
}
