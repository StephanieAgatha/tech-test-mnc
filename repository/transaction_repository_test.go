package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"mnc-test/model"
	"reflect"
	"testing"
	"time"
)

func TestMakePayment_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTransactionRepository(db)

	tx := &model.Transaction{
		CustomerID:    1,
		MerchantID:    2,
		BankAccountID: 3,
		Amount:        1000,
		MerchantName:  "Test Merchant",
	}

	mock.ExpectBegin()
	mock.ExpectExec("insert into transactions (.+) values (.+)").WithArgs(tx.CustomerID, tx.MerchantID, tx.BankAccountID, tx.Amount).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update BankAccounts SET balance = balance - \\$1 where id = \\$2 and customer_id = \\$3").WithArgs(tx.Amount, tx.BankAccountID, tx.CustomerID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update MerchantBalances SET balance = balance \\+ \\$1 where merchant_id = \\$2").WithArgs(tx.Amount, tx.MerchantID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT name FROM Merchants WHERE id = \\$1").WithArgs(tx.MerchantID).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(tx.MerchantName))
	mock.ExpectCommit()

	err = r.MakePayment(tx)
	if err != nil {
		t.Errorf("error was not expected while making payment: %s", err)
	}
}

func TestMakePayment_InsertTransaction_FailureEntryTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTransactionRepository(db)

	tx := &model.Transaction{
		CustomerID:    1,
		MerchantID:    2,
		BankAccountID: 3,
		Amount:        1000,
		MerchantName:  "Test Merchant",
	}

	mock.ExpectBegin()
	mock.ExpectExec("insert into transactions (.+) values (.+)").WithArgs(tx.CustomerID, tx.MerchantID, tx.BankAccountID, tx.Amount).WillReturnError(fmt.Errorf("insert error")) // simulate an error
	mock.ExpectRollback()

	err = r.MakePayment(tx)
	if err == nil {
		t.Errorf("an error was expected when inserting transaction")
	} else if err.Error() != "insert error" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

// test pengurangan saldo failure / gagal
func TestMakePayment_DeductBalanceFailure_DeductBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTransactionRepository(db)

	tx := &model.Transaction{
		CustomerID:    1,
		MerchantID:    2,
		BankAccountID: 3,
		Amount:        1000,
		MerchantName:  "Test Merchant",
	}

	mock.ExpectBegin()
	mock.ExpectExec("insert into transactions (.+) values (.+)").WithArgs(tx.CustomerID, tx.MerchantID, tx.BankAccountID, tx.Amount).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update BankAccounts SET balance = balance - \\$1 where id = \\$2 and customer_id = \\$3").WithArgs(tx.Amount, tx.BankAccountID, tx.CustomerID).WillReturnError(fmt.Errorf("deduct error")) // simulate an error
	mock.ExpectRollback()

	err = r.MakePayment(tx)
	if err == nil {
		t.Errorf("an error was expected when deducting balance")
	} else if err.Error() != "deduct error" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

func TestMakePayment_AddBalanceFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTransactionRepository(db)

	tx := &model.Transaction{
		CustomerID:    1,
		MerchantID:    2,
		BankAccountID: 3,
		Amount:        1000,
		MerchantName:  "Test Merchant",
	}

	mock.ExpectBegin()
	mock.ExpectExec("insert into transactions (.+) values (.+)").WithArgs(tx.CustomerID, tx.MerchantID, tx.BankAccountID, tx.Amount).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update BankAccounts SET balance = balance - \\$1 where id = \\$2 and customer_id = \\$3").WithArgs(tx.Amount, tx.BankAccountID, tx.CustomerID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update MerchantBalances SET balance = balance \\+ \\$1 where merchant_id = \\$2").WithArgs(tx.Amount, tx.MerchantID).WillReturnError(fmt.Errorf("add balance error")) // simulate an error
	mock.ExpectRollback()

	err = r.MakePayment(tx)
	if err == nil {
		t.Errorf("an error was expected when adding balance")
	} else if err.Error() != "add balance error" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

func TestMakePayment_GetMerchantNameFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTransactionRepository(db)

	tx := &model.Transaction{
		CustomerID:    1,
		MerchantID:    2,
		BankAccountID: 3,
		Amount:        1000,
		MerchantName:  "Test Merchant",
	}

	mock.ExpectBegin()
	mock.ExpectExec("insert into transactions (.+) values (.+)").WithArgs(tx.CustomerID, tx.MerchantID, tx.BankAccountID, tx.Amount).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update BankAccounts SET balance = balance - \\$1 where id = \\$2 and customer_id = \\$3").WithArgs(tx.Amount, tx.BankAccountID, tx.CustomerID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update MerchantBalances SET balance = balance \\+ \\$1 where merchant_id = \\$2").WithArgs(tx.Amount, tx.MerchantID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT name FROM Merchants WHERE id = \\$1").WithArgs(tx.MerchantID).WillReturnError(fmt.Errorf("get merchant name error")) // simulate an error
	mock.ExpectRollback()

	err = r.MakePayment(tx)
	if err == nil {
		t.Errorf("an error was expected when getting merchant name")
	} else if err.Error() != "get merchant name error" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

func TestMakePayment_CommitFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTransactionRepository(db)

	tx := &model.Transaction{
		CustomerID:    1,
		MerchantID:    2,
		BankAccountID: 3,
		Amount:        1000,
		MerchantName:  "Test Merchant",
	}

	mock.ExpectBegin()
	mock.ExpectExec("insert into transactions (.+) values (.+)").WithArgs(tx.CustomerID, tx.MerchantID, tx.BankAccountID, tx.Amount).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update BankAccounts SET balance = balance - \\$1 where id = \\$2 and customer_id = \\$3").WithArgs(tx.Amount, tx.BankAccountID, tx.CustomerID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("update MerchantBalances SET balance = balance \\+ \\$1 where merchant_id = \\$2").WithArgs(tx.Amount, tx.MerchantID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT name FROM Merchants WHERE id = \\$1").WithArgs(tx.MerchantID).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(tx.MerchantName))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("commit error")) // simulate an error
	mock.ExpectRollback()

	err = r.MakePayment(tx)
	if err == nil {
		t.Errorf("an error was expected when committing transaction")
	} else if err.Error() != "commit error" {
		t.Errorf("unexpected error returned: %s", err)
	}
}

// test pada func get tx
func TestGetCustomerTransactionByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTransactionRepository(db)

	custID := 1
	tx := &model.Transaction{
		ID:            1,
		CustomerID:    custID,
		MerchantID:    2,
		BankAccountID: 3,
		Amount:        1000,
		MerchantName:  "Test Merchant",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "customer_id", "merchant_id", "bank_account_id", "amount", "created_at", "updated_at", "name"}).
		AddRow(tx.ID, tx.CustomerID, tx.MerchantID, tx.BankAccountID, tx.Amount, tx.CreatedAt, tx.UpdatedAt, tx.MerchantName)

	mock.ExpectQuery(`^SELECT t\.\*, m\.name FROM Transactions t JOIN Merchants m ON t\.merchant_id = m\.id WHERE t\.customer_id = \$1 ORDER BY t\.created_at DESC$`).
		WithArgs(custID).
		WillReturnRows(rows)

	result, err := r.GetCustomerTransactionByID(custID)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err)
	}

	if len(result) != 1 {
		t.Errorf("expected one transaction, got %d", len(result))
	} else if !reflect.DeepEqual(result[0], *tx) {
		t.Errorf("expected %v, got %v", *tx, result[0])
	}
}

func TestGetCustomerTransactionByID_FailureRowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTransactionRepository(db)

	custID := 1
	tx := &model.Transaction{
		ID:            1,
		CustomerID:    custID,
		MerchantID:    2,
		BankAccountID: 3,
		Amount:        1000,
		MerchantName:  "Test Merchant",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Note that we're only returning 7 columns instead of 8
	rows := sqlmock.NewRows([]string{"id", "customer_id", "merchant_id", "bank_account_id", "amount", "created_at", "updated_at"}).
		AddRow(tx.ID, tx.CustomerID, tx.MerchantID, tx.BankAccountID, tx.Amount, tx.CreatedAt, tx.UpdatedAt)

	mock.ExpectQuery(`^SELECT t\.\*, m\.name FROM Transactions t JOIN Merchants m ON t\.merchant_id = m\.id WHERE t\.customer_id = \$1 ORDER BY t\.created_at DESC$`).
		WithArgs(custID).
		WillReturnRows(rows)

	result, err := r.GetCustomerTransactionByID(custID)
	if err == nil {
		t.Errorf("expected an error, but got none")
	}

	if result != nil {
		t.Errorf("expected result to be nil, but got %v", result)
	}
}
