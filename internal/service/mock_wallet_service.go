package service

import (
	"github.com/google/uuid"
)

// MockWalletService - мок-реализация WalletService
type MockWalletService struct {
	ProcessOperationFn func(walletID uuid.UUID, operationType string, amount int64) error
	GetBalanceFn       func(walletID uuid.UUID) (int64, error)
}

func (m *MockWalletService) ProcessOperation(walletID uuid.UUID, operationType string, amount int64) error {
	if m.ProcessOperationFn != nil {
		return m.ProcessOperationFn(walletID, operationType, amount)
	}
	return nil
}

func (m *MockWalletService) GetBalance(walletID uuid.UUID) (int64, error) {
	if m.GetBalanceFn != nil {
		return m.GetBalanceFn(walletID)
	}
	return 0, nil
}
