package usecase

import (
	"mnc-test/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransferRepository struct {
	mock.Mock
}

func (m *MockTransferRepository) MakeTransferAccNumbToAccNumb(transactionID, senderAccountNumber, receiverAccountNumber string, amount int) error {
	args := m.Called(transactionID, senderAccountNumber, receiverAccountNumber, amount)
	return args.Error(0)
}

func (m *MockTransferRepository) GetIncomingMoney(customerId int) ([]model.TransferHistoryIncome, error) {
	args := m.Called(customerId)
	return args.Get(0).([]model.TransferHistoryIncome), args.Error(1)
}

func (m *MockTransferRepository) GetOutcomeMoney(customerId int) ([]model.TransferHistoryOutcome, error) {
	args := m.Called(customerId)
	return args.Get(0).([]model.TransferHistoryOutcome), args.Error(1)
}

func TestTransferUsecase_MakeTransferAccNumbToAccNumb_Success(t *testing.T) {
	mockRepo := new(MockTransferRepository)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("MakeTransferAccNumbToAccNumb", mock.Anything, "senderAccNum", "receiverAccNum", 1000).Return(nil)

		u := NewTransferUsecase(mockRepo)
		transactionID, err := u.MakeTransferAccNumbToAccNumb("senderAccNum", "receiverAccNum", 1000)

		assert.NoError(t, err)
		assert.NotEmpty(t, transactionID)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransferUsecase_GetIncomingMoney_Success(t *testing.T) {
	mockRepo := new(MockTransferRepository)
	mockIncomingMoney := []model.TransferHistoryIncome{
		{CustomerID: 1, Amount: 1000},
		{CustomerID: 1, Amount: 2000},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetIncomingMoney", 1).Return(mockIncomingMoney, nil)

		u := NewTransferUsecase(mockRepo)
		result, err := u.GetIncomingMoney(1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransferUsecase_GetOutcomeMoney_Success(t *testing.T) {
	mockRepo := new(MockTransferRepository)
	mockOutcomeMoney := []model.TransferHistoryOutcome{
		{CustomerID: 1, Amount: 1000},
		{CustomerID: 1, Amount: 2000},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetOutcomeMoney", 1).Return(mockOutcomeMoney, nil)

		u := NewTransferUsecase(mockRepo)
		result, err := u.GetOutcomeMoney(1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
