package usecase

import (
	"mnc-test/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMerchantRepository struct {
	mock.Mock
}

func (m *MockMerchantRepository) CreateNewMerchant(merchant model.Merchant) error {
	args := m.Called(merchant)
	return args.Error(0)
}

func (m *MockMerchantRepository) FindAllMerchant() ([]model.Merchant, error) {
	args := m.Called()
	return args.Get(0).([]model.Merchant), args.Error(1)
}

func TestMerchantUsecase_CreateNewMerchant(t *testing.T) {
	mockRepo := new(MockMerchantRepository)
	mockMerchant := model.Merchant{Name: "Test Merchant"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CreateNewMerchant", mockMerchant).Return(nil)

		u := NewMerchantUsecase(mockRepo)
		err := u.CreateNewMerchant(mockMerchant)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

}

func TestMerchantUsecase_FindAllMerchant(t *testing.T) {
	mockRepo := new(MockMerchantRepository)
	mockMerchants := []model.Merchant{
		{Name: "Test Merchant 1"},
		{Name: "Test Merchant 2"},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("FindAllMerchant").Return(mockMerchants, nil)

		u := NewMerchantUsecase(mockRepo)
		result, err := u.FindAllMerchant()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})

}
