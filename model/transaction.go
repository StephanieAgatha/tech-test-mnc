package model

import "time"

type Transaction struct {
	ID            int       `json:"id"`
	CustomerID    int       `json:"customer_id"`
	MerchantID    int       `json:"merchant_id"`
	BankAccountID int       `json:"bank_account_id"`
	Amount        int       `json:"amount"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CustomerEmail string    `json:"customer_email,omitempty"`
	MerchantName  string    `json:"merchant_name,omitempty"`
}

type TransactionResponse struct {
	CustomerID    int    `json:"customer_id,omitempty"`
	TransactionID string `json:"transaction_id"`
	Amount        int    `json:"amount"`
	CreatedAt     string `json:"created_at"`
	MerchantName  string `json:"merchant_name"`
}
