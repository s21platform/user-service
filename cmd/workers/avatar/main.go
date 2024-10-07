package main

import (
	"log"

	"github.com/s21platform/user-service/internal/config"
	newavatar "github.com/s21platform/user-service/internal/repository/kafka/consumer/new_avatar"
	"github.com/s21platform/user-service/internal/repository/postgres"
)

func main() {
	cfg := config.MustLoad()

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	avatarUpdater, err := newavatar.New(cfg, dbRepo)
	if err != nil {
		log.Fatalf("error create consumer: %v", err)
	}

	avatarUpdater.Listen()
}
