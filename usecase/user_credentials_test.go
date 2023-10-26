package usecase

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"mnc-test/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserCredentialRepository struct {
	mock.Mock
}

func (m *MockUserCredentialRepository) Login(userCred model.UserCredentials) (string, error) {
	args := m.Called(userCred)
	return args.String(0), args.Error(1)
}

func (m *MockUserCredentialRepository) FindUserEMail(email string) (userCred model.UserCredentials, err error) {
	args := m.Called(email)
	return args.Get(0).(model.UserCredentials), args.Error(1)
}

func (m *MockUserCredentialRepository) Register(userCred model.UserCredentials) error {
	args := m.Called(userCred)
	return args.Error(0)
}

func TestUserCredentialUsecase_Register(t *testing.T) {
	mockRepo := new(MockUserCredentialRepository)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //isi jika ada pw nya
		DB:       0,
	})

	t.Run("success", func(t *testing.T) {
		user := model.UserCredentials{
			Name:     "test",
			Email:    "test@test.com",
			Password: "Password1!",
		}

		mockRepo.On("Register", mock.MatchedBy(func(user model.UserCredentials) bool {
			return user.Name == "test" && user.Email == "test@test.com"
		})).Return(nil)

		u := NewUserCredentialUsecase(mockRepo, client)
		err := u.Register(user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	//test when we got empty email
	t.Run("error-emptyemail", func(t *testing.T) {
		user := model.UserCredentials{
			Name:     "test",
			Email:    "",
			Password: "Password1!",
		}

		u := NewUserCredentialUsecase(mockRepo, client)
		err := u.Register(user)

		assert.Error(t, err)
	})

	//password is empty
	t.Run("error-emptypass", func(t *testing.T) {
		user := model.UserCredentials{
			Name:     "test",
			Email:    "awd@gmail.com",
			Password: "",
		}

		u := NewUserCredentialUsecase(mockRepo, client)
		err := u.Register(user)

		assert.Error(t, err)
	})

	t.Run("error-shortpassword", func(t *testing.T) {
		user := model.UserCredentials{
			Name:     "test",
			Email:    "test@test.com",
			Password: "Pwd1!",
		}

		u := NewUserCredentialUsecase(mockRepo, client)
		err := u.Register(user)

		assert.Error(t, err)
		assert.Equal(t, "Password must contain at least six number", err.Error())
	})

	t.Run("error-no uppercase letter", func(t *testing.T) {
		user := model.UserCredentials{
			Name:     "test",
			Email:    "test@test.com",
			Password: "password1!",
		}

		u := NewUserCredentialUsecase(mockRepo, client)
		err := u.Register(user)

		assert.Error(t, err)
		assert.Equal(t, "Password must contain at least one uppercase letter", err.Error())
	})

	t.Run("error-no special character", func(t *testing.T) {
		user := model.UserCredentials{
			Name:     "test",
			Email:    "test@test.com",
			Password: "Password1",
		}

		u := NewUserCredentialUsecase(mockRepo, client)
		err := u.Register(user)

		assert.Error(t, err)
		assert.Equal(t, "Password must contain at least one special character", err.Error())
	})

	t.Run("error-no number", func(t *testing.T) {
		user := model.UserCredentials{
			Name:     "test",
			Email:    "test@test.com",
			Password: "Password!",
		}

		u := NewUserCredentialUsecase(mockRepo, client)
		err := u.Register(user)

		assert.Error(t, err)
		assert.Equal(t, "Password must contain at least one number", err.Error())
	})
}

// finduseremail test
func TestUserCredentialUsecase_FindUserEmail(t *testing.T) {
	mockRepo := new(MockUserCredentialRepository)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //isi jika ada pw nya
		DB:       0,
	})

	t.Run("success", func(t *testing.T) {
		email := "test@helowww.com"

		user := model.UserCredentials{
			Email:    email,
			Password: "Password1!",
		}

		mockRepo.On("FindUserEMail", mock.MatchedBy(func(email string) bool {
			return email == "test@helowww.com"
		})).Return(user, nil)

		u := NewUserCredentialUsecase(mockRepo, client)
		foundUser, err := u.FindUserEMail(email)

		assert.NoError(t, err)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, user.Password, foundUser.Password)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - no email", func(t *testing.T) {
		email := ""

		u := NewUserCredentialUsecase(mockRepo, client)
		foundUser, err := u.FindUserEMail(email)

		assert.Error(t, err)
		assert.Equal(t, model.UserCredentials{}, foundUser)
		assert.Equal(t, "Email is required", err.Error())
	})

	t.Run("error - user not found", func(t *testing.T) {
		email := "abra@cyaa.com"

		mockRepo.On("FindUserEMail", mock.MatchedBy(func(email string) bool {
			return email == "abra@cyaa.com"
		})).Return(model.UserCredentials{}, errors.New("user not found"))

		u := NewUserCredentialUsecase(mockRepo, client)
		foundUser, err := u.FindUserEMail(email)

		assert.Error(t, err)
		assert.Equal(t, model.UserCredentials{}, foundUser)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
