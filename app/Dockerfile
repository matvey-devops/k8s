# Build stage
FROM golang:1.21-alpine AS builder

# Установка необходимых пакетов для сборки
RUN apk add --no-cache git ca-certificates

# Создание рабочей директории
WORKDIR /build

# Копирование go mod файлов
COPY go.mod ./

# Загрузка зависимостей (если будут)
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

# Установка необходимых пакетов
RUN apk --no-cache add ca-certificates tzdata

# Создание непривилегированного пользователя
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

# Копирование бинарного файла из builder stage
COPY --from=builder /build/main /app/main

# Установка прав на исполнение
RUN chmod +x /app/main

# Смена пользователя
USER appuser

# Рабочая директория
WORKDIR /app

# Порт по умолчанию
EXPOSE 8080

# Запуск приложения
CMD ["./main"] 