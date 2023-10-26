package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// makepayment func test
func TestMakeTransferAccNumbToAccNumb_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	transactionID := "466d6803-1fb4-4cca-a630-bf1322c36bb0"
	senderAccountNumber := "12481257"
	receiverAccountNumber := "12371246"
	amount := 10000

	mock.ExpectBegin()
	mock.ExpectQuery("^select balance from bankaccounts where account_number = \\$1 for update").WithArgs(senderAccountNumber).WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(20000))
	mock.ExpectExec("^update bankaccounts set balance = balance - \\$1 where account_number = \\$2").WithArgs(amount, senderAccountNumber).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^update bankaccounts set balance = balance \\+ \\$1 where account_number = \\$2").WithArgs(amount, receiverAccountNumber).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^insert into transfer_history \\(transaction_id, sender_account_number, receiver_account_number, amount\\) values \\(\\$1, \\$2, \\$3, \\$4\\)").WithArgs(transactionID, senderAccountNumber, receiverAccountNumber, amount).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	r := &transferRepository{db: db}
	err = r.MakeTransferAccNumbToAccNumb(transactionID, senderAccountNumber, receiverAccountNumber, amount)
	if err != nil {
		t.Errorf("an error was not expected when making transfer: %s", err)
	}
}

func TestMakeTransferAccNumbToAccNumb_InsufficientBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	transactionID := "466d6803-1fb4-4cca-a630-bf1322c36bb0"
	senderAccountNumber := "12481257"
	receiverAccountNumber := "12371246"
	amount := 10000

	mock.ExpectBegin()
	mock.ExpectQuery("^select balance from bankaccounts where account_number = \\$1 for update").WithArgs(senderAccountNumber).WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(5000)) // the balance is less than the amount
	mock.ExpectRollback()

	r := &transferRepository{db: db}
	err = r.MakeTransferAccNumbToAccNumb(transactionID, senderAccountNumber, receiverAccountNumber, amount)
	if err == nil {
		t.Errorf("an error was expected when making transfer due to insufficient balance")
	} else if err.Error() != "Insufficient balance" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

func TestMakeTransferAccNumbToAccNumb_QueryBalanceSenderError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	transactionID := "466d6803-1fb4-4cca-a630-bf1322c36bb0"
	senderAccountNumber := "12481257"
	receiverAccountNumber := "12371246"
	amount := 10000

	mock.ExpectBegin()
	mock.ExpectQuery("^select balance from bankaccounts where account_number = \\$1 for update").WithArgs(senderAccountNumber).WillReturnError(errors.New("database error")) // simulate a database error
	mock.ExpectRollback()

	r := &transferRepository{db: db}
	err = r.MakeTransferAccNumbToAccNumb(transactionID, senderAccountNumber, receiverAccountNumber, amount)
	if err == nil {
		t.Errorf("an error was expected when making transfer due to database error")
	} else if err.Error() != "database error" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

func TestMakeTransferAccNumbToAccNumb_UpdateSenderBalanceError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	transactionID := "466d6803-1fb4-4cca-a630-bf1322c36bb0"
	senderAccountNumber := "12481257"
	receiverAccountNumber := "12371246"
	amount := 10000

	mock.ExpectBegin()
	mock.ExpectQuery("^select balance from bankaccounts where account_number = \\$1 for update").WithArgs(senderAccountNumber).WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(20000))
	mock.ExpectExec("^update bankaccounts set balance = balance - \\$1 where account_number = \\$2").WithArgs(amount, senderAccountNumber).WillReturnError(errors.New("database error")) // simulate a database error
	mock.ExpectRollback()

	r := &transferRepository{db: db}
	err = r.MakeTransferAccNumbToAccNumb(transactionID, senderAccountNumber, receiverAccountNumber, amount)
	if err == nil {
		t.Errorf("an error was expected when making transfer due to database error")
	} else if err.Error() != "database error" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

// when updating the receiver's balance is failsss
func TestMakeTransferAccNumbToAccNumb_UpdateReceiverBalanceError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	transactionID := "466d6803-1fb4-4cca-a630-bf1322c36bb0"
	senderAccountNumber := "12481257"
	receiverAccountNumber := "12371246"
	amount := 10000

	mock.ExpectBegin()
	mock.ExpectQuery("^select balance from bankaccounts where account_number = \\$1 for update").WithArgs(senderAccountNumber).WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(20000))
	mock.ExpectExec("^update bankaccounts set balance = balance - \\$1 where account_number = \\$2").WithArgs(amount, senderAccountNumber).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^update bankaccounts set balance = balance \\+ \\$1 where account_number = \\$2").WithArgs(amount, receiverAccountNumber).WillReturnError(errors.New("database error")) // simulate a database error
	mock.ExpectRollback()

	r := &transferRepository{db: db}
	err = r.MakeTransferAccNumbToAccNumb(transactionID, senderAccountNumber, receiverAccountNumber, amount)
	if err == nil {
		t.Errorf("an error was expected when making transfer due to database error")
	} else if err.Error() != "database error" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

// when inserting into transfer_history table fails
func TestMakeTransferAccNumbToAccNumb_InsertTransferHistoryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	transactionID := "466d6803-1fb4-4cca-a630-bf1322c36bb0"
	senderAccountNumber := "12481257"
	receiverAccountNumber := "12371246"
	amount := 10000

	mock.ExpectBegin()
	mock.ExpectQuery("^select balance from bankaccounts where account_number = \\$1 for update").WithArgs(senderAccountNumber).WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(20000))
	mock.ExpectExec("^update bankaccounts set balance = balance - \\$1 where account_number = \\$2").WithArgs(amount, senderAccountNumber).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^update bankaccounts set balance = balance \\+ \\$1 where account_number = \\$2").WithArgs(amount, receiverAccountNumber).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^insert into transfer_history \\(transaction_id, sender_account_number, receiver_account_number, amount\\) values \\(\\$1, \\$2, \\$3, \\$4\\)").WithArgs(transactionID, senderAccountNumber, receiverAccountNumber, amount).WillReturnError(errors.New("database error")) // simulate a database error
	mock.ExpectRollback()

	r := &transferRepository{db: db}
	err = r.MakeTransferAccNumbToAccNumb(transactionID, senderAccountNumber, receiverAccountNumber, amount)
	if err == nil {
		t.Errorf("an error was expected when making transfer due to database error")
	} else if err.Error() != "database error" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

// incomingmoney func test
func TestGetIncomingMoney_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	//parse timestamp
	timestamp, err := time.Parse("2006-01-02 15:04:05.999999", "2023-10-26 10:49:51.391191")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing timestamp", err)
	}

	//mock rows
	rows := sqlmock.NewRows([]string{"id", "transaction_id", "sender_account_number", "receiver_account_number", "amount", "transfer_timestamp", "sender_bank_name", "sender_bank_id", "receiver_bank_name", "receiver_bank_id"}).
		AddRow(15, "466d6803-1fb4-4cca-a630-bf1322c36bb0", "12481257", "12371246", 10000, timestamp, "BCA", 1, "BRI", 2)

	//query goes here
	mock.ExpectQuery("^select account_number from bankaccounts where customer_id = \\$1").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"account_number"}).AddRow("12481257"))
	mock.ExpectQuery("^select (.+) from transfer_history th join bankaccounts ba1 on th.sender_account_number = ba1.account_number join banks b1 on ba1.bank_id = b1.id join bankaccounts ba2 on th.receiver_account_number = ba2.account_number join banks b2 on ba2.bank_id = b2.id where ba1.customer_id = \\$1 OR ba2.customer_id = \\$1").WithArgs(5).WillReturnRows(rows)

	r := &transferRepository{db: db}
	result, err := r.GetIncomingMoney(5)
	if err != nil {
		t.Errorf("error was not expected while getting incoming money: %s", err)
	}

	if len(result) > 0 {
		assert.Equal(t, 15, result[0].ID)
		assert.Equal(t, "466d6803-1fb4-4cca-a630-bf1322c36bb0", result[0].TransactionID)
		assert.Equal(t, "12481257", result[0].SenderAccountNumber)
		assert.Equal(t, "12371246", result[0].ReceiverAccountNumber)
		assert.Equal(t, 10000, result[0].Amount)
		assert.Equal(t, timestamp, result[0].TransferTimeStamp)
		assert.Equal(t, "BCA", result[0].SenderBankName)
		assert.Equal(t, 1, result[0].SenderBankId)
		assert.Equal(t, "BRI", result[0].ReceiverBankName)
		assert.Equal(t, 2, result[0].ReceiverBankId)
	}
}

func TestGetIncomingMoney_FailureFetchAccountNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("^select account_number from bankaccounts where customer_id = \\$1").
		WithArgs(5).
		WillReturnError(errors.New("some error"))

	r := &transferRepository{db: db}
	result, err := r.GetIncomingMoney(5)
	if err == nil {
		t.Errorf("an error was expected while getting incoming money, got nil")
	}

	if result != nil {
		t.Errorf("result was not expected to have any value, got: %v", result)
	}
}

func TestGetIncomingMoney_FailureScanRowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timestamp, err := time.Parse("2006-01-02 15:04:05.999999", "2023-10-26 10:49:51.391191")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing timestamp", err)
	}

	// Add an extra column to cause a scan error
	rows := sqlmock.NewRows([]string{"id", "transaction_id", "sender_account_number", "receiver_account_number", "amount", "transfer_timestamp", "sender_bank_name", "sender_bank_id", "receiver_bank_name", "receiver_bank_id", "blablaa"}).
		AddRow(15, "466d6803-1fb4-4cca-a630-bf1322c36bb0", "12481257", "12371246", 10000, timestamp, "BCA", 1, "BRI", 2, "cyaaa")

	mock.ExpectQuery("^select account_number from bankaccounts where customer_id = \\$1").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"account_number"}).AddRow("12481257"))
	mock.ExpectQuery("^select (.+) from transfer_history th join bankaccounts ba1 on th.sender_account_number = ba1.account_number join banks b1 on ba1.bank_id = b1.id join bankaccounts ba2 on th.receiver_account_number = ba2.account_number join banks b2 on ba2.bank_id = b2.id where ba1.customer_id = \\$1 OR ba2.customer_id = \\$1").WithArgs(5).WillReturnRows(rows)

	r := &transferRepository{db: db}
	result, err := r.GetIncomingMoney(5)
	if err == nil {
		t.Errorf("an error was expected while scanning rows, got nil")
	}

	if result != nil {
		t.Errorf("result was not expected to have any value, got: %v", result)
	}
}

func TestGetIncomingMoney_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Add an extra column to cause a scan error
	rows := sqlmock.NewRows([]string{"id", "transaction_id", "sender_account_number", "receiver_account_number", "amount", "transfer_timestamp", "sender_bank_name", "sender_bank_id", "receiver_bank_name", "receiver_bank_id"}).
		AddRow(15, "466d6803-1fb4-4cca-a630-bf1322c36bb0", "12481257", "12371246", 10000, time.Now(), "BCA", 1, "BRI", 2).
		RowError(0, errors.New("some error")) // This will cause rows.Next() to return an error

	mock.ExpectQuery("^select account_number from bankaccounts where customer_id = \\$1").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"account_number"}).AddRow("12481257"))
	mock.ExpectQuery("^select (.+) from transfer_history th join bankaccounts ba1 on th.sender_account_number = ba1.account_number join banks b1 on ba1.bank_id = b1.id join bankaccounts ba2 on th.receiver_account_number = ba2.account_number join banks b2 on ba2.bank_id = b2.id where ba1.customer_id = \\$1 OR ba2.customer_id = \\$1").WithArgs(5).WillReturnRows(rows)

	r := &transferRepository{db: db}
	result, err := r.GetIncomingMoney(5)
	if err == nil {
		t.Errorf("an error was expected while getting incoming money, got nil")
	}

	if result != nil {
		t.Errorf("result was not expected to have any value, got: %v", result)
	}
}

func TestGetOutcomeMoney_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	timestamp, err := time.Parse("2006-01-02 15:04:05.999999", "2023-10-26 10:49:51.391191")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when parsing timestamp", err)
	}

	rows := sqlmock.NewRows([]string{"id", "transaction_id", "sender_account_number", "receiver_account_number", "amount", "transfer_timestamp", "sender_bank_name", "sender_bank_id", "receiver_bank_name", "receiver_bank_id"}).
		AddRow(15, "466d6803-1fb4-4cca-a630-bf1322c36bb0", "12481257", "12371246", 10000, timestamp, "BCA", 1, "BRI", 2)

	mock.ExpectQuery("^select (.+) from bankaccounts where customer_id = \\$1").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"account_number"}).AddRow("12481257"))
	mock.ExpectQuery("^select (.+) from transfer_history th join bankaccounts ba1 on th.sender_account_number = ba1.account_number join banks b1 on ba1.bank_id = b1.id join bankaccounts ba2 on th.receiver_account_number = ba2.account_number join banks b2 on ba2.bank_id = b2.id where ba1.customer_id = \\$1").WithArgs(5).WillReturnRows(rows)

	r := &transferRepository{db: db}
	result, err := r.GetOutcomeMoney(5)
	if err != nil {
		t.Errorf("error was not expected while getting outcome money: %s", err)
	}

	if len(result) > 0 {
		assert.Equal(t, 15, result[0].ID)
		assert.Equal(t, "466d6803-1fb4-4cca-a630-bf1322c36bb0", result[0].TransactionID)
		assert.Equal(t, "12481257", result[0].SenderAccountNumber)
		assert.Equal(t, "12371246", result[0].ReceiverAccountNumber)
		assert.Equal(t, 10000, result[0].Amount)
		assert.Equal(t, timestamp, result[0].TransferTimeStamp)
		assert.Equal(t, "BCA", result[0].SenderBankName)
		assert.Equal(t, 1, result[0].SenderBankId)
		assert.Equal(t, "BRI", result[0].ReceiverBankName)
		assert.Equal(t, 2, result[0].ReceiverBankId)
	} else {
		t.Errorf("result was expected to have at least one entry, got %d", len(result))
	}
}
