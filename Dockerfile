# Используем базовый образ на основе Alpine для сборки и выполнения
FROM golang:1.23-alpine

# Установка зависимостей
RUN apk update && apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование файлов go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Установка зависимостей
RUN go mod download

# Копирование всего исходного кода
COPY . .

# Копирование конфигурационного файла
COPY .env .env

# Установка переменных окружения
ENV CONFIG_PATH=/app/.env

# Указание команды для выполнения
CMD [ "go", "run", "./cmd/main" ]
