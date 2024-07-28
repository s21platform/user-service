package main

import (
	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/repository/postgres"
)

func main() {
	cfg := config.MustLoad()

	db := postgres.New(cfg)
	defer db.Close()
}
