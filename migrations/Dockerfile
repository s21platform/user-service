# Используем официальный образ Golang как базовый
FROM golang:1.22-alpine

# Установим необходимые зависимости
RUN apk add --no-cache bash gcc musl-dev postgresql-dev

# Установим Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Создадим рабочую директорию
WORKDIR /app

# Скопируем файлы миграций
COPY . /app/db/migrations

# Скопируем entrypoint скрипт
COPY entrypoint.sh /app/entrypoint.sh

# Сделаем entrypoint скрипт исполняемым
RUN chmod +x /app/entrypoint.sh

# Установим entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]