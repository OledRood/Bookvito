# Этап 1: Сборка приложения (Builder)
# Используем образ с нужной версией Go (1.23), как указано в go.mod
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые для сборки утилиты
RUN apk add --no-cache git build-base

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей и загружаем их
# Это кэширует зависимости, и они не будут скачиваться каждый раз, если go.mod/go.sum не менялись
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код проекта
COPY . .

# Убедимся что go.sum актуален
RUN go mod tidy

# Собираем бинарный файл приложения
# CGO_ENABLED=0 отключает Cgo, что делает бинарник статически скомпонованным
# Флаги -w -s уменьшают размер бинарного файла
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /app/main ./cmd/api

# Этап 2: Финальный образ (Final)
# Используем минимальный образ alpine для уменьшения размера финального образа
FROM alpine:latest

# Добавляем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем пользователя без прав root для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/main .
RUN chown appuser:appgroup /app/main
USER appuser

EXPOSE 8080
CMD ["./main"]