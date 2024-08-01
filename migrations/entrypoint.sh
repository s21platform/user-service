#!/bin/bash

ls /app/db/migrations

# Применим миграции
goose -dir /app/db/migrations postgres "user=$USER_SERVICE_POSTGRES_USER password=$USER_SERVICE_POSTGRES_PASSWORD dbname=$USER_SERVICE_POSTGRES_DB host=$USER_SERVICE_POSTGRES_HOST port=$USER_SERVICE_POSTGRES_PORT sslmode=disable" up
