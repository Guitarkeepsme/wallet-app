package service_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-app/internal/handler"
	"wallet-app/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	zeroAmount       = 0
	zeroAmountError  = "amount must be greater than zero"
	unknownOperation = "unknown operation type"
)

// Mock для репозитория
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Deposit(walletID uuid.UUID, amount int64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockRepository) Withdraw(walletID uuid.UUID, amount int64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockRepository) GetBalance(walletID uuid.UUID) (int64, error) {
	args := m.Called(walletID)
	return args.Get(0).(int64), args.Error(1)
}

func TestProcessOperation_Deposit(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.New(mockRepo)

	walletID := uuid.New()
	amount := int64(100)
	mockRepo.On("Deposit", walletID, amount).Return(nil)

	err := service.ProcessOperation(walletID, "DEPOSIT", amount)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Deposit", walletID, amount)
}

func TestProcessOperation_Deposit_InvalidAmount(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.New(mockRepo)

	walletID := uuid.New()
	amount := int64(-100) // Неправильная сумма

	err := service.ProcessOperation(walletID, "DEPOSIT", amount)

	assert.Error(t, err)
	assert.Equal(t, zeroAmountError, err.Error())
	mockRepo.AssertNotCalled(t, "Deposit")
}

func TestProcessOperation_Deposit_InvalidDataType(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.New(mockRepo)

	// Симулируем некорректный тип данных
	walletID := uuid.New()
	var amount interface{} = ";sakldfj"

	// Конвертация в int64 вызовет панику или ошибку; важно обработать это на уровне вызова сервиса
	assert.Panics(t, func() {
		_ = service.ProcessOperation(walletID, "DEPOSIT", amount.(int64))
	})

	mockRepo.AssertNotCalled(t, "Deposit")
}

func TestHandler_InvalidAmount(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.New(mockRepo)
	h := handler.New(service)

	server := httptest.NewServer(http.HandlerFunc(h.HandleWalletOperation))
	defer server.Close()

	// Некорректное значение суммы
	reqBody := []byte(`{
		"walletID": "d290f1ee-6c54-4b01-90e6-d701748f0851",
		"operationType": "DEPOSIT",
		"amount": ";sakldfj"
	}`)
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestProcessOperation_Withdraw_InsufficientAmount(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.New(mockRepo)

	walletID := uuid.New()
	amount := int64(-50)

	err := service.ProcessOperation(walletID, "WITHDRAW", amount)

	assert.Error(t, err)
	assert.Equal(t, zeroAmountError, err.Error())
}

func TestProcessOperation_Withdraw_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.New(mockRepo)

	walletID := uuid.New()
	amount := int64(100)
	mockRepo.On("Withdraw", walletID, amount).Return(nil)

	err := service.ProcessOperation(walletID, "WITHDRAW", amount)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Withdraw", walletID, amount)
}

func TestGetBalance_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.New(mockRepo)

	walletID := uuid.New()
	expectedBalance := int64(1000)
	mockRepo.On("GetBalance", walletID).Return(expectedBalance, nil)

	balance, err := service.GetBalance(walletID)

	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	mockRepo.AssertCalled(t, "GetBalance", walletID)
}

func TestGetBalance_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.New(mockRepo)
	expectedErr := errors.New("wallet not found")

	walletID := uuid.New()
	mockRepo.On("GetBalance", walletID).Return(int64(0), expectedErr)

	balance, err := service.GetBalance(walletID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, int64(zeroAmount), balance)
	mockRepo.AssertCalled(t, "GetBalance", walletID)
}

func TestProcessOperation_InvalidOperation(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.New(mockRepo)

	walletID := uuid.New()
	amount := int64(100)

	err := service.ProcessOperation(walletID, "INVALID_TYPE", amount)

	assert.Error(t, err)
	assert.Equal(t, "unknown operation type, INVALID_TYPE", err.Error())
	mockRepo.AssertNotCalled(t, "Deposit")
	mockRepo.AssertNotCalled(t, "Withdraw")
}
