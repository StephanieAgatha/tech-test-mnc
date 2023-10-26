package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"mnc-test/model"
	"testing"
	"time"
)

func TestCreateNewMerchant(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new merchantRepository
	repo := NewMerchantRepository(db)

	// Mock merchant
	merchant := model.Merchant{
		Name: "Test Merchant",
	}

	mock.ExpectExec("^insert into merchants \\(name\\) values \\(\\$1\\)$").
		WithArgs(merchant.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateNewMerchant(merchant)

	require.NoError(t, err)

	// Assert that all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNewMerchantFailure(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	// Create a new MerchantRepository
	repo := NewMerchantRepository(db)

	// Mock merchant
	merchant := model.Merchant{
		Name: "Test Merchant",
	}

	mock.ExpectExec("^insert into merchants \\(name\\) values \\(\\$1\\)$").
		WithArgs(merchant.Name).
		WillReturnError(fmt.Errorf("some database error"))

	// Call the method
	err = repo.CreateNewMerchant(merchant)

	// Assert that an error was returned
	require.Error(t, err)
	require.Equal(t, "Failed to exec query some database error", err.Error())

	// Assert that all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllMerchant(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	// Create a new repo
	repo := NewMerchantRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(1, "Test Merchant", time.Now(), time.Now())
	mock.ExpectQuery("^select \\* from merchants$").
		WillReturnRows(rows)

	// Call the method
	merchants, err := repo.FindAllMerchant()

	// Assert that no error was returned
	require.NoError(t, err)

	// Assert that one merchant was returned
	require.Len(t, merchants, 1)

	// Assert that the returned merchant is the one we expected
	require.Equal(t, "Test Merchant", merchants[0].Name)

	// Assert that all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllMerchantFailure(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMerchantRepository(db)

	mock.ExpectQuery("^select \\* from merchants$").
		WillReturnError(fmt.Errorf("some database error"))

	// Call the method
	_, err = repo.FindAllMerchant()

	// Assert that an error was returned
	require.Error(t, err)

	// Assert that the error is due to a failed query execution
	require.Equal(t, "Failed to find merchants some database error", err.Error())

	// Assert that all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllMerchantScanFailure(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new MerchantRepository
	repo := NewMerchantRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow("invalid", "invalid", "invalid", "invalid")
	mock.ExpectQuery("^select \\* from merchants$").
		WillReturnRows(rows)

	// Call the method
	_, err = repo.FindAllMerchant()

	// Assert that an error was returned
	require.Error(t, err)

	// Assert that the error is due to a failed retrieval of merchants
	require.Contains(t, err.Error(), "Failed to retrieve merchants")

	// Assert that all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}
