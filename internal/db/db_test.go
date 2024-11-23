package db_test

import (
	"database/sql"
	"testing"
	"wallet-app/internal/db"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	// Инициализируем in-memory базу данных
	database, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer database.Close()

	err = db.NewDB(database)
	assert.NoError(t, err)

	// Проверяем, что таблица создана
	var tableName string
	err = database.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='wallets';").Scan(&tableName)
	assert.NoError(t, err)
	assert.Equal(t, "wallets", tableName)
}
