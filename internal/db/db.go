package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Можно вынести перенести формирование строки подключения и протестировать её

func ConnectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	return sql.Open("postgres", connStr)
}

// MigrateDB - функция для выполнения миграции
// сделать мок-коннект в бд и проверить, правильный ли запрос улетает

func NewDB(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS wallets (
			wallet_id UUID PRIMARY KEY,
			balance BIGINT NOT NULL CHECK (balance >= 0)
		);
	`
	_, err := db.Exec(query)
	return err
}
