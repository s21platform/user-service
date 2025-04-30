package main

import (
	"context"
	"log"

	kafkalib "github.com/s21platform/kafka-lib"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/databus/new_avatar"
	"github.com/s21platform/user-service/internal/repository/postgres"
)

func main() {
	cfg := config.MustLoad()

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "user", cfg.Platform.Env)
	if err != nil {
		log.Println("failed to connect graphite: ", err)
	}

	ctx := context.WithValue(context.Background(), config.KeyMetrics, metrics)

	consumerCfg := kafkalib.DefaultConsumerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.SetNewAvatar, "user")
	newAvatarConsumer, err := kafkalib.NewConsumer(consumerCfg, metrics)
	if err != nil {
		log.Println("error create consumer: ", err)
	}

	newAvatarHandler := new_avatar.New(dbRepo)

	newAvatarConsumer.RegisterHandler(ctx, newAvatarHandler.Handler)

	<-ctx.Done()
}
