package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"wallet-app/internal/db"
	"wallet-app/internal/handler"
	"wallet-app/internal/repository"
	"wallet-app/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file: ", err)
	}

	// Подключаем базу данных
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	database.SetMaxOpenConns(100)                // Максимальное количество соединений
	database.SetMaxIdleConns(50)                 // Количество соединений в режиме ожидания
	database.SetConnMaxLifetime(5 * time.Minute) // Время жизни соединения

	// Закрываем базу данных только после завершения работы
	defer database.Close()

	// Выполняем миграции
	err = db.NewDB(database)
	if err != nil {
		log.Fatal("Ошибка выполнения миграций:", err)
	}

	// Инициализируем зависимости
	repo := repository.New(database)
	service := service.New(repo)
	handler := handler.New(service)

	// Настраиваем маршруты
	router := chi.NewRouter()
	router.Post("/api/v1/wallet", handler.HandleWalletOperation)
	router.Get("/api/v1/wallets/{walletID}", handler.GetWalletBalance)

	// Отображаем все маршруты
	chi.Walk(router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

	// Запускаем сервер
	port := os.Getenv("APP_PORT")
	log.Printf("Server is running on port %s", port)

	http.ListenAndServe(":"+port, router)
}
