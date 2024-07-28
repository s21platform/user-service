package main

import (
	"github.com/s21platform/user-service/internal/config"
	"log"
)

func main() {
	cfg := config.MustLoad()
	log.Println("config", cfg)
}
