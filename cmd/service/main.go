package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/s21platform/metrics-lib/pkg"
	user "github.com/s21platform/user-proto/user-proto"

	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/infra"
	"github.com/s21platform/user-service/internal/repository/kafka/producer/friends_register"
	"github.com/s21platform/user-service/internal/repository/postgres"
	"github.com/s21platform/user-service/internal/service"
)

func main() {
	cfg := config.MustLoad()

	db := postgres.New(cfg)
	defer db.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "user", cfg.Platform.Env)
	if err != nil {
		log.Fatalf("failed to create metrics object: %v", err)
	}

	userFriendsRegister := friends_register.New(cfg)

	server := service.New(db, userFriendsRegister)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.UnaryInterceptor,
			infra.MetricsInterceptor(metrics),
		),
	)
	user.RegisterUserServiceServer(s, server)

	log.Println("starting server", cfg.Service.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
