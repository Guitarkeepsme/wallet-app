# Используем базовый образ Golang
FROM golang:1.23.2 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod download

# Собираем исполняемый файл
RUN go build -o wallet-app ./cmd/main.go

# Инициализируем финальный образ
FROM debian:bookworm-slim

WORKDIR /app

# Копируем скомпилирвоанное приложеное
COPY --from=builder /app/wallet-app .
COPY .env .

# Указываем порт приложения
EXPOSE ${APP_PORT}

# Запускаем приложение
CMD ["./wallet-app"]
