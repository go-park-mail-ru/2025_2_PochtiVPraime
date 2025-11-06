# Stage 1: Сборка приложения
FROM golang:1.25.1-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Stage 2: Запуск приложения
FROM alpine:latest

# Устанавливаем зависимости времени выполнения
RUN apk --no-cache add ca-certificates


WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /app/main .

# Экспонируем порт (замените на ваш порт)
EXPOSE 8080

# Команда запуска
CMD ["./main"]