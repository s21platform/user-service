package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"

	kafkalib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	optoinhub "github.com/s21platform/user-service/internal/clients/optionhub"
	"github.com/s21platform/user-service/internal/config"
	api "github.com/s21platform/user-service/internal/generated"
	"github.com/s21platform/user-service/internal/infra"
	"github.com/s21platform/user-service/internal/pkg/tx"
	"github.com/s21platform/user-service/internal/repository/postgres"
	"github.com/s21platform/user-service/internal/rest"
	"github.com/s21platform/user-service/internal/service"
	"github.com/s21platform/user-service/pkg/user"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger.Error(fmt.Sprintf("cannnot listen port; error: %s", err))
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}

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

	UserPostCreatedProduserConfig := kafkalib.DefaultProducerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.UserPostCreated)
	UserPostCreatedProduser := kafkalib.NewProducer(UserPostCreatedProduserConfig)

	server := service.New(db, producerNewFriendRegister, optionhubClient, prcUserCreated, UserPostCreatedProduser)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.Logger(logger),
			infra.UnaryInterceptor,
			infra.MetricsInterceptor(metrics),
			tx.TxMiddleWire(db),
		),
	)
	user.RegisterUserServiceServer(s, server)

	h := rest.New()
	r := chi.NewRouter()
	api.HandlerFromMux(h, r)
	httpServer := &http.Server{
		Handler: r,
	}

	m := cmux.New(lis)

	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	go func() {
		if err := s.Serve(grpcL); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	go func() {
		if err := httpServer.Serve(httpL); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	log.Println("starting server", cfg.Service.Port)
	if err := m.Serve(); err != nil {
		logger.Error(fmt.Sprintf("cannot start service; error: %s", err))
		log.Fatalf("Cannnot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
