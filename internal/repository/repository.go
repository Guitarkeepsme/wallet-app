package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

const (
	ErrInsufficientFunds = "insufficient funds"
	ErrWalletNotFound    = "wallet not found"
)

type WalletRepository interface {
	GetBalance(walletID uuid.UUID) (int64, error)
	Deposit(walletID uuid.UUID, amount int64) error
	Withdraw(walletID uuid.UUID, amount int64) error
}

// Реализация интерфейса для работы с базой данных
type walletRepositoryImpl struct {
	db *sql.DB
}

func New(db *sql.DB) WalletRepository {
	return &walletRepositoryImpl{db: db}
}

// Кладём деньги, обращаясь к базе данных
func (r *walletRepositoryImpl) Deposit(walletID uuid.UUID, amount int64) error {
	_, err := r.db.Exec(`
	    INSERT INTO wallets (wallet_id, balance)
		VALUES ($1, $2)
		ON CONFLICT (wallet_id) DO UPDATE
		SET balance = wallets.balance + $2
	`, walletID, amount)

	return err
}

// Снимаем деньги
// Может, здесь нужно amount float64?
func (r *walletRepositoryImpl) Withdraw(walletID uuid.UUID, amount int64) error {
	// Здесь нам нужно сохранить результат транзакции, потому что мы сделаем несколько запросов
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var balance int64
	err = tx.QueryRow(`
		SELECT balance 
		FROM wallets 
		WHERE wallet_id = $1
		FOR UPDATE`, walletID).Scan(&balance)
	// FOR UPDATE

	if err != nil {
		tx.Rollback()
		return err
	}

	if balance < amount {
		tx.Rollback()
		return errors.New(ErrInsufficientFunds)
	}

	_, err = tx.Exec(`
        UPDATE wallets 
        SET balance = balance - $1 
        WHERE wallet_id = $2`, amount, walletID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Завершаем транзакцию
	return tx.Commit()
}

// Получаем баланс
func (r *walletRepositoryImpl) GetBalance(walletID uuid.UUID) (int64, error) {
	var balance int64
	err := r.db.QueryRow(`
        SELECT balance 
        FROM wallets 
        WHERE wallet_id = $1`, walletID).Scan(&balance)

	if err == sql.ErrNoRows {
		return 0, errors.New(ErrWalletNotFound)
	}

	return balance, err
}
