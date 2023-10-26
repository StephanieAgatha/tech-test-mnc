package repository

import (
	"database/sql"
	"fmt"
	"mnc-test/model"
)

type UserCredential interface {
	Register(userCred model.UserCredentials) error
	Login(userCred model.UserCredentials) (string, error)
	FindUserEMail(email string) (userCred model.UserCredentials, err error)
}

type userCredential struct {
	db *sql.DB
}

func (u userCredential) FindUserEMail(email string) (userCred model.UserCredentials, err error) {
	//TODO implement me

	query := "select id,email,phone_number,password from Customers where email = $1"

	if err = u.db.QueryRow(query, email).Scan(&userCred.ID, &userCred.Email, &userCred.PhoneNumber, &userCred.Password); err != nil {
		if err == sql.ErrNoRows {
			return model.UserCredentials{}, fmt.Errorf("Invalid Credential")
		}
		return model.UserCredentials{}, fmt.Errorf("Failed to run query %v", err.Error())
	}
	return userCred, nil
}

func (u userCredential) Register(userCred model.UserCredentials) error {
	//TODO implement me
	query := "insert into Customers (name,email,phone_number,password) values ($1, $2, $3, $4)"

	_, err := u.db.Exec(query, userCred.Name, userCred.Email, userCred.PhoneNumber, userCred.Password)
	if err != nil {
		return fmt.Errorf("Failed to exec query %v", err.Error())
	}
	return nil
}

func (u userCredential) Login(userCred model.UserCredentials) (string, error) {
	//TODO implement me

	var hashedPass string
	query := "select password from Customers where email = $1 "
	err := u.db.QueryRow(query, userCred.Email).Scan(&hashedPass)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Invalid Credentials %v", err.Error())
		}
		return "", fmt.Errorf("Failed to exec query")
	}
	return hashedPass, nil
}

func NewUserCredentials(db *sql.DB) UserCredential {
	return &userCredential{
		db: db,
	}
}
