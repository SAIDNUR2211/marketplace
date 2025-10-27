# --- Этап 1: Сборка приложения ---
FROM golang:1.24.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь остальной исходный код
COPY . .

# Собираем статически скомпилированный бинарный файл для Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go


# --- Этап 2: Создание финального легковесного образа ---
FROM alpine:latest

WORKDIR /app

# Копируем только скомпилированный бинарный файл из этапа сборки
COPY --from=builder /main .

# Копируем необходимые для работы файлы конфигурации
# Dockerfile должен находиться в корне, чтобы этот путь был верным
COPY internal/configs/configs.json ./internal/configs/configs.json

# Открываем порт, на котором работает приложение
EXPOSE 7577

# Команда для запуска приложения при старте контейнера
CMD ["./main"]