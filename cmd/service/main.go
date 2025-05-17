package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	kafkalib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"
	"github.com/s21platform/user-service/pkg/user"

	optoinhub "github.com/s21platform/user-service/internal/clients/optionhub"
	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/infra"
	"github.com/s21platform/user-service/internal/repository/postgres"
	"github.com/s21platform/user-service/internal/service"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	db := postgres.New(cfg)
	defer db.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "user", cfg.Platform.Env)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create metrics object: %v", err))
		log.Fatalf("failed to create metrics object: %v", err)
	}

	producerCfg := kafkalib.DefaultProducerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.FriendsRegister)
	producerNewFriendRegister := kafkalib.NewProducer(producerCfg)
	optionhubClient := optoinhub.MustConnect(cfg)

	prcUserCreatedCfg := kafkalib.DefaultProducerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.UserCreated)
	prcUserCreated := kafkalib.NewProducer(prcUserCreatedCfg)

	server := service.New(db, producerNewFriendRegister, optionhubClient, prcUserCreated)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.Logger(logger),
			infra.UnaryInterceptor,
			infra.MetricsInterceptor(metrics),
		),
	)
	user.RegisterUserServiceServer(s, server)

	log.Println("starting server", cfg.Service.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger.Error(fmt.Sprintf("cannnot listen port; error: %s", err))
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}
	if err := s.Serve(lis); err != nil {
		logger.Error(fmt.Sprintf("cannot start service; error: %s", err))
		log.Fatalf("Cannnot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
