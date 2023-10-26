package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"mnc-test/model"
	"testing"
)

func TestRegister(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new UserCredential repository
	repo := NewUserCredentials(db)

	// Mock user credentials
	userCred := model.UserCredentials{
		Name:        "Test User",
		Email:       "test@example.com",
		PhoneNumber: "085156810985",
		Password:    "password",
	}

	mock.ExpectExec("^insert into Customers \\(name,email,phone_number,password\\) values \\(\\$1, \\$2, \\$3, \\$4\\)$").
		WithArgs(userCred.Name, userCred.Email, userCred.PhoneNumber, userCred.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Register(userCred)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRegisterFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserCredentials(db)

	// Mock user credentials
	userCred := model.UserCredentials{
		Name:        "Test User",
		Email:       "test@example.com",
		PhoneNumber: "1234567890",
		Password:    "password",
	}

	mock.ExpectExec("^insert into Customers \\(name,email,phone_number,password\\) values \\(\\$1, \\$2, \\$3, \\$4\\)$").
		WithArgs(userCred.Name, userCred.Email, userCred.PhoneNumber, userCred.Password).
		WillReturnError(fmt.Errorf("some database error"))

	err = repo.Register(userCred)
	require.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	// Create a new UserCredential repository
	repo := NewUserCredentials(db)

	// Mock user credentials
	userCred := model.UserCredentials{
		Email:    "test@example.com",
		Password: "password",
	}

	// Mock hashed password
	hashedPass := "hashedpassword"

	// Expect a select query
	mock.ExpectQuery("^select password from Customers where email = \\$1$").
		WithArgs(userCred.Email).
		WillReturnRows(sqlmock.NewRows([]string{"password"}).
			AddRow(hashedPass))

	password, err := repo.Login(userCred)
	require.NoError(t, err)

	require.Equal(t, hashedPass, password)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserCredentials(db)

	// Mock user credentials
	userCred := model.UserCredentials{
		Email:    "test@example.com",
		Password: "password",
	}

	mock.ExpectQuery("^select password from Customers where email = \\$1$").
		WithArgs(userCred.Email).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.Login(userCred)

	require.Error(t, err)
	require.Equal(t, "Invalid Credentials sql: no rows in result set", err.Error())

	// Reset expectations
	mock.ExpectationsWereMet()

	mock.ExpectQuery("^select password from Customers where email = \\$1$").
		WithArgs(userCred.Email).
		WillReturnError(fmt.Errorf("some database error"))

	_, err = repo.Login(userCred)
	require.Error(t, err)

	// Assert that the error is due to a failed query execution
	require.Equal(t, "Failed to exec query", err.Error())
	// Assert that all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindUserEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserCredentials(db)

	// Mock user credentials
	userCred := model.UserCredentials{
		ID:          "1",
		Email:       "test@example.com",
		PhoneNumber: "1234567890",
		Password:    "password",
	}

	mock.ExpectQuery("^select id,email,phone_number,password from Customers where email = \\$1$").
		WithArgs(userCred.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "phone_number", "password"}).
			AddRow(userCred.ID, userCred.Email, userCred.PhoneNumber, userCred.Password))

	// Call the FindUserEmail method
	result, err := repo.FindUserEMail(userCred.Email)

	require.NoError(t, err)

	require.Equal(t, userCred, result)
	// Assert that all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}
