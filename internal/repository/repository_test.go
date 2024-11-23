package repository_test

import (
	"testing"
	"wallet-app/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWalletRepository_Deposit(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.New(db)

	walletID := uuid.New()
	mock.ExpectExec("INSERT INTO wallets").
		WithArgs(walletID, 1000).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Deposit(walletID, 1000)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWalletRepository_Withdraw(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.New(db)

	walletID := uuid.New()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT balance").
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(5000))
	mock.ExpectExec("UPDATE wallets").
		WithArgs(1000, walletID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Withdraw(walletID, 1000)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWalletRepository_GetBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.New(db)

	walletID := uuid.New()
	mock.ExpectQuery("SELECT balance").
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(5000))

	balance, err := repo.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(5000), balance)
}
