package usecase

import (
	"mnc-test/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) MakePayment(tx *model.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetCustomerTransactionByID(custID int) ([]model.Transaction, error) {
	args := m.Called(custID)
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func TestTransactionUsecase_MakePayment_Success(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	mockTransaction := &model.Transaction{CustomerID: 1, MerchantID: 1, Amount: 1000}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("MakePayment", mockTransaction).Return(nil)

		u := NewTransactionUsecase(mockRepo, nil)
		err := u.MakePayment(mockTransaction)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionUsecase_GetCustomerTransaction_Success(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	mockTransactions := []model.Transaction{
		{CustomerID: 1, MerchantID: 1, Amount: 1000},
		{CustomerID: 1, MerchantID: 2, Amount: 2000},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetCustomerTransactionByID", 1).Return(mockTransactions, nil)

		u := NewTransactionUsecase(mockRepo, nil)
		result, err := u.GetCustomerTransaction(1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
