package usecase

import (
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
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
	redisClient    *redis.Client
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
	pasetoToken := helper.GeneratePaseto(userCred.Email, symetricKey)

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
	token := userCred.Token

	if token == "" {
		return errors.New("invalid token")
	}

	// blacklist them
	u.mu.Lock()
	u.tokenBlacklist[token] = true
	u.mu.Unlock()

	return nil
}

func NewUserCredentialUsecase(usercredRepo repository.UserCredential, redisClient *redis.Client) UserCredentialUsecase {
	return &userCredentialUsecase{
		usercredRepo: usercredRepo,
		redisClient:  redisClient,
	}
}
