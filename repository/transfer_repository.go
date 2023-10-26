package repository

import (
	"context"
	"database/sql"
	"errors"
	"mnc-test/model"
)

type TransferRepository interface {
	MakeTransferAccNumbToAccNumb(transactionID string, senderAccountNumber string, receiverAccountNumber string, amount int) error
	GetIncomingMoney(customerId int) ([]model.TransferHistoryIncome, error)
	GetOutcomeMoney(customerId int) ([]model.TransferHistoryOutcome, error)
}

type transferRepository struct {
	db *sql.DB
}

func (t transferRepository) MakeTransferAccNumbToAccNumb(transactionID string, senderAccountNumber string, receiverAccountNumber string, amount int) error {
	// Start a new transaction
	tx, err := t.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	// Rollback if anything goes wrong
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	// Check if the sender has enough balance
	var senderBalance int
	err = tx.QueryRow("select balance from bankaccounts where account_number = $1 for update", senderAccountNumber).Scan(&senderBalance)
	if err != nil {
		tx.Rollback()
		return err
	}

	if senderBalance < amount {
		tx.Rollback()
		return errors.New("Insufficient balance")
	}

	// update the sender's balance
	_, err = tx.Exec("update bankaccounts set balance = balance - $1 where account_number = $2", amount, senderAccountNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	// update the receiver's balance
	_, err = tx.Exec("update bankaccounts set balance = balance + $1 where account_number = $2", amount, receiverAccountNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert the transfer details into the transfer_history table
	_, err = tx.Exec("insert into transfer_history (transfer_id, sender_account_number, receiver_account_number, amount) values ($1, $2, $3, $4)", transactionID, senderAccountNumber, receiverAccountNumber, amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (t transferRepository) GetIncomingMoney(customerId int) ([]model.TransferHistoryIncome, error) {
	var accountNumber string

	// Fetch the account number
	err := t.db.QueryRow("select account_number from bankaccounts where customer_id = $1", customerId).Scan(&accountNumber)
	if err != nil {
		return nil, err
	}

	rows, err := t.db.Query(`
	select th.id, th.transfer_id, th.sender_account_number, th.receiver_account_number, th.amount, th.transfer_timestamp, 
	b1.name as sender_bank_name, b1.id as sender_bank_id,
	b2.name as receiver_bank_name, b2.id as receiver_bank_id 
	from transfer_history th
	join bankaccounts ba1 on th.sender_account_number = ba1.account_number
	join banks b1 on ba1.bank_id = b1.id
	join bankaccounts ba2 on th.receiver_account_number = ba2.account_number
	join banks b2 on ba2.bank_id = b2.id
	where ba1.customer_id = $1 OR ba2.customer_id = $1`, customerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transfers := []model.TransferHistoryIncome{}

	for rows.Next() {
		var transfer model.TransferHistoryIncome
		err = rows.Scan(&transfer.ID, &transfer.TransferID, &transfer.SenderAccountNumber, &transfer.ReceiverAccountNumber, &transfer.Amount, &transfer.TransferTimeStamp, &transfer.SenderBankName, &transfer.SenderBankId, &transfer.ReceiverBankName, &transfer.ReceiverBankId)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transfers, nil
}

func (t transferRepository) GetOutcomeMoney(customerId int) ([]model.TransferHistoryOutcome, error) {
	var accountNumber string

	// Fetch the account number
	err := t.db.QueryRow("select account_number from bankaccounts where customer_id = $1", customerId).Scan(&accountNumber)
	if err != nil {
		return nil, err
	}

	rows, err := t.db.Query(`
	select th.id, th.transfer_id, th.sender_account_number, th.receiver_account_number, th.amount, th.transfer_timestamp, 
	b1.name as sender_bank_name, b1.id as sender_bank_id,
	b2.name as receiver_bank_name, b2.id as receiver_bank_id 
	from transfer_history th
	join bankaccounts ba1 on th.sender_account_number = ba1.account_number
	join banks b1 on ba1.bank_id = b1.id
	join bankaccounts ba2 on th.receiver_account_number = ba2.account_number
	join banks b2 on ba2.bank_id = b2.id
	where ba1.customer_id = $1`, customerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transfers := []model.TransferHistoryOutcome{}

	for rows.Next() {
		var transfer model.TransferHistoryOutcome
		err = rows.Scan(&transfer.ID, &transfer.TransactionID, &transfer.SenderAccountNumber, &transfer.ReceiverAccountNumber, &transfer.Amount, &transfer.TransferTimeStamp, &transfer.SenderBankName, &transfer.SenderBankId, &transfer.ReceiverBankName, &transfer.ReceiverBankId)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transfers, nil
}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{
		db: db,
	}
}

//func (t transferRepository) MakeTransferPhoneNumbToPhoneNumb(transactionID string, senderPhoneNumber string, receiverPhoneNumber string, amount int) error {
//	// Start a new transaction
//	tx, err := t.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
//	if err != nil {
//		return err
//	}
//
//	// Rollback if anything goes wrong
//	defer func() {
//		if p := recover(); p != nil {
//			tx.Rollback()
//		}
//	}()
//
//	var senderCustomerId, receiverCustomerId int
//	var senderAccountNumber, receiverAccountNumber string
//
//	// Fetch the sender's customer ID
//	err = tx.QueryRow("select id from customers where phone_number = $1", senderPhoneNumber).Scan(&senderCustomerId)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	// Fetch the sender's account number
//	err = tx.QueryRow("select account_number from bankaccounts where customer_id = $1", senderCustomerId).Scan(&senderAccountNumber)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	// Fetch the receiver's customer ID
//	err = tx.QueryRow("select id from customers where phone_number = $1", receiverPhoneNumber).Scan(&receiverCustomerId)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	// Fetch the receiver's account number
//	err = tx.QueryRow("select account_number from bankaccounts where customer_id = $1", receiverCustomerId).Scan(&receiverAccountNumber)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	// Check if the sender has enough balance
//	var senderBalance int
//	err = tx.QueryRow("select balance from bankaccounts where account_number = $1 FOR update", senderAccountNumber).Scan(&senderBalance)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	if senderBalance < amount {
//		tx.Rollback()
//		return errors.New("Insufficient balance")
//	}
//
//	// update the sender's balance
//	_, err = tx.Exec("update bankaccounts SET balance = balance - $1 where account_number = $2", amount, senderAccountNumber)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	// update the receiver's balance
//	_, err = tx.Exec("update bankaccounts SET balance = balance + $1 where account_number = $2", amount, receiverAccountNumber)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	// Insert the transfer details into the transfer_history table
//	_, err = tx.Exec("INSERT INTO transfer_history (transaction_id, sender_account_number, receiver_account_number, amount) VALUES ($1, $2, $3, $4)", transactionID, senderAccountNumber, receiverAccountNumber, amount)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	err = tx.Commit()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
