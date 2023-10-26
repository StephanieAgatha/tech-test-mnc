package repository

import (
	"database/sql"
	"mnc-test/model"
)

type TransactionRepository interface {
	MakePayment(tx *model.Transaction) error
	GetCustomerTransactionByID(custID int) ([]model.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func (t *transactionRepository) MakePayment(tx *model.Transaction) error {
	// Start a new transaction
	dbTx, err := t.db.Begin()
	if err != nil {
		return err
	}

	// Step 1: Create a new entry in the transactions table
	_, err = dbTx.Exec(`insert into transactions (customer_id, merchant_id, bank_account_id, amount) values ($1, $2, $3, $4)`,
		tx.CustomerID, tx.MerchantID, tx.BankAccountID, tx.Amount)
	if err != nil {
		dbTx.Rollback()
		return err
	}

	// Step 2: Deduct the amount from the bank account's balance
	_, err = dbTx.Exec(`update BankAccounts SET balance = balance - $1 where id = $2 and customer_id = $3`,
		tx.Amount, tx.BankAccountID, tx.CustomerID)
	if err != nil {
		dbTx.Rollback()
		return err
	}

	// Step 3: Add the amount to the merchant's balance
	_, err = dbTx.Exec(`update MerchantBalances SET balance = balance + $1 where merchant_id = $2`,
		tx.Amount, tx.MerchantID)
	if err != nil {
		dbTx.Rollback()
		return err
	}

	// Step 4: Get the merchant name
	err = dbTx.QueryRow(`SELECT name FROM Merchants WHERE id = $1`, tx.MerchantID).Scan(&tx.MerchantName)
	if err != nil {
		dbTx.Rollback()
		return err
	}

	// Commit the transaction
	err = dbTx.Commit()
	if err != nil {
		dbTx.Rollback()
		return err
	}

	return nil
}

func (t *transactionRepository) GetCustomerTransactionByID(custID int) ([]model.Transaction, error) {
	rows, err := t.db.Query(`
		SELECT t.*, m.name 
		FROM Transactions t 
		JOIN Merchants m ON t.merchant_id = m.id 
		WHERE t.customer_id = $1 
		ORDER BY t.created_at DESC
	`, custID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var tx model.Transaction
		err = rows.Scan(&tx.ID, &tx.CustomerID, &tx.MerchantID, &tx.BankAccountID, &tx.Amount, &tx.CreatedAt, &tx.UpdatedAt, &tx.MerchantName) // Added MerchantName
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}
