package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
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
	ctx := logger_lib.NewContext(context.Background(), logger)

	db := postgres.New(cfg)
	defer db.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "user", cfg.Platform.Env)
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to create metrics object: %v", err))
		log.Fatalf("failed to create metrics object: %v", err)
	}

	producerCfg := kafkalib.DefaultProducerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.FriendsRegister)
	producerNewFriendRegister := kafkalib.NewProducer(producerCfg)
	optionhubClient := optoinhub.MustConnect(cfg)

	prcUserCreatedCfg := kafkalib.DefaultProducerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.UserCreated)
	prcUserCreated := kafkalib.NewProducer(prcUserCreatedCfg)

	UserPostCreatedProducerConfig := kafkalib.DefaultProducerConfig(cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.UserPostCreated)
	UserPostCreatedProducer := kafkalib.NewProducer(UserPostCreatedProducerConfig)

	server := service.New(db, producerNewFriendRegister, prcUserCreated, UserPostCreatedProducer)

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.Logger(logger),
			infra.UnaryInterceptor,
			infra.MetricsInterceptor(metrics),
			tx.TxMiddleWire(db),
		),
	)
	user.RegisterUserServiceServer(grpcSrv, server)

	handler := rest.New(db, optionhubClient)
	router := chi.NewRouter()
	router.Use(infra.AuthRequest)
	router.Use(infra.LoggerRequest(logger))

	api.HandlerFromMux(handler, router)
	httpServer := &http.Server{
		Handler: router,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("cannnot listen port; error: %s", err))
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}

	m := cmux.New(lis)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	g, _ := errgroup.WithContext(context.Background())

	logger_lib.Info(ctx, "starting server")

	g.Go(func() error {
		if err := grpcSrv.Serve(grpcListener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("gRPC server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := httpServer.Serve(httpListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := m.Serve(); err != nil {
			return fmt.Errorf("cannot start service: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("server error: %v", err))
		log.Fatalf("Server error: %v", err)
	}
}
