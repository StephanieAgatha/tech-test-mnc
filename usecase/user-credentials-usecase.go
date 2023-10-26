package usecase

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"mnc-test/model"
	"mnc-test/repository"
	"mnc-test/util/helper"
	"os"
	"sync"
)

type UserCredentialUsecase interface {
	Register(userCred model.UserCredentials) error
	Login(userCred model.UserCredentials) (string, error)
	Logout(userCred model.UserCredentials) error
	FindUserEMail(email string) (userCred model.UserCredentials, err error)
}

type userCredentialUsecase struct {
	usercredRepo   repository.UserCredential
	tokenBlacklist map[string]bool
	mu             sync.Mutex // for concurrent map writes
	log            *zap.Logger
}

func (u *userCredentialUsecase) Register(userCred model.UserCredentials) error {
	//TODO implement me

	//generate uuid for user id
	userCred.ID = helper.GenerateUUID()

	if userCred.Email == "" {
		return fmt.Errorf("Username is required")
	}

	if userCred.Password == "" {
		return fmt.Errorf("Password is required")
	}

	//is email alr valid?
	if err := helper.IsEmailValid(userCred); err != nil {
		return err
	}

	/*
		password requirement
	*/
	if len(userCred.Password) < 6 {
		return fmt.Errorf("Password must contain at least six number")
	}
	if !helper.PasswordContainsUppercase(userCred.Password) {
		return fmt.Errorf("Password must contain at least one uppercase letter")
	}

	if !helper.PasswordContainsSpecialChar(userCred.Password) {
		return fmt.Errorf("Password must contain at least one special character")
	}

	if !helper.PasswordConstainsOneNumber(userCred.Password) {
		return fmt.Errorf("Password must contain at least one number")
	}

	//generate password in here
	hashedPass, err := helper.HashPassword(userCred.Password)
	if err != nil {
		return err
	}

	userCred.Password = hashedPass

	if err = u.usercredRepo.Register(userCred); err != nil {
		return err
	}

	//log goes hereeee
	if u.log != nil {
		u.log.Info("New Customer Has Been Created",
			zap.String("Custormer Name", userCred.Name),
			zap.String("Customer Email", userCred.Email))
	} else {
		fmt.Println("Logger is not initialized")
	}

	return nil
}

func (u *userCredentialUsecase) Login(userCred model.UserCredentials) (string, error) {
	//TODO implement me

	if userCred.Email == "" {
		return "", fmt.Errorf("Email is required")
	} else if userCred.Password == "" {
		return "", fmt.Errorf("Password is required")
	}

	userHashedPass, err := u.usercredRepo.Login(userCred)
	if err != nil {

	}
	//compare password
	if err = helper.ComparePassword(userHashedPass, userCred.Password); err != nil {
		return "", fmt.Errorf("Invalid Password")
	}

	//generate paseto or jwt in here
	symetricKey := os.Getenv("PASETO_SECRET")
	//jwtsecret := os.Getenv("JWT_SECRET")
	pasetoToken := helper.GeneratePaseto(userCred.Email, symetricKey)

	if u.log != nil {
		u.log.Info("Customer has been logged in",
			zap.String("Customer Email", userCred.Email))
	} else {
		fmt.Println("Logger is not initialized")
	}

	return pasetoToken, nil
}

func (u *userCredentialUsecase) FindUserEMail(email string) (userCred model.UserCredentials, err error) {
	//TODO implement me

	if email == "" {
		return model.UserCredentials{}, fmt.Errorf("Email is required")
	}

	user, err := u.usercredRepo.FindUserEMail(email)
	if err != nil {
		return model.UserCredentials{}, err
	}

	return user, nil
}

func (u *userCredentialUsecase) Logout(userCred model.UserCredentials) error {
	// Extract the token from the user credentials
	// This will depend on how you have structured your UserCredentials model
	token := userCred.Token

	if token == "" {
		return errors.New("invalid token")
	}

	// Add the token to the blacklist
	u.mu.Lock()
	u.tokenBlacklist[token] = true
	u.mu.Unlock()

	return nil
}

func NewUserCredentialUsecase(uc repository.UserCredential, log *zap.Logger) UserCredentialUsecase {
	return &userCredentialUsecase{
		usercredRepo:   uc,
		tokenBlacklist: make(map[string]bool),
		mu:             sync.Mutex{},
		log:            log,
	}
}
