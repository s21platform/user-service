# Используем базовый образ Ubuntu
FROM ubuntu:latest

# Обновляем пакеты и устанавливаем необходимые зависимости
RUN apt-get update && \
    apt-get install -y \
    bash \
    ca-certificates \
    curl \
    git \
    gcc \
    g++ \
    libc-dev \
    postgresql-client && \
    rm -rf /var/lib/apt/lists/*

# Устанавливаем golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

# Копирование миграций в контейнер
WORKDIR /app
COPY . ./migrations

# Установка entrypoint
ENTRYPOINT ["/usr/local/bin/migrate"]