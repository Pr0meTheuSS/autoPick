# Используем официальный образ Golang для сборки
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod tidy

# Собираем статически с поддержкой CGO для совместимости с Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o content-service ./cmd/api

# Финальный контейнер на основе `alpine`
FROM alpine:latest

WORKDIR /app

# Копируем собранный бинарник
COPY --from=builder /app/content-service .

# Делаем бинарник исполняемым (если вдруг проблема с правами)
RUN chmod +x ./content-service

# Копируем .env файл
COPY .env .

# Открываем порт
EXPOSE 9090

# Запускаем сервис
CMD ["/app/content-service"]
