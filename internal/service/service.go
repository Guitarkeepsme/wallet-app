package service

import (
	"errors"
	"fmt"
	"sync"
	"wallet-app/internal/repository"

	"github.com/google/uuid"
)

const (
	zeroAmount       = 0
	zeroAmountError  = "amount must be greater than zero"
	unknownOperation = "unknown operation type"
)

var ErrWalletNotFound = errors.New("wallet not found")
var ErrInvalidOperation = errors.New("invalid amount")

type WalletService interface {
	ProcessOperation(walletID uuid.UUID, operationType string, amount int64) error
	GetBalance(walletID uuid.UUID) (int64, error)
}

type walletServiceImpl struct {
	repo repository.WalletRepository
}

func New(repo repository.WalletRepository) WalletService {
	return &walletServiceImpl{repo: repo}
}

// Для предотвращения гонки данных используем мьютексы
var walletLocks sync.Map

// Проводим операцию над балансом кошелька
func (s *walletServiceImpl) ProcessOperation(walletID uuid.UUID, operationType string, amount int64) error {
	// Получаем или создаём мьютекс для данного walletID
	lock, _ := walletLocks.LoadOrStore(walletID, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)

	// Блокируем доступ к кошельку
	mutex.Lock()
	defer mutex.Unlock()

	// Проверяем, что баланс положительный
	if amount <= zeroAmount {
		return fmt.Errorf("%s", zeroAmountError)
	}
	// Смотрим, какую именно операцию получили
	switch operationType {
	case "DEPOSIT":
		return s.repo.Deposit(walletID, amount)
	case "WITHDRAW":
		return s.repo.Withdraw(walletID, amount)
	default:
		return fmt.Errorf("%s, %s", unknownOperation, operationType)
	}
}

// Проверяем баланс
func (s *walletServiceImpl) GetBalance(walletID uuid.UUID) (int64, error) {
	return s.repo.GetBalance(walletID)
}
