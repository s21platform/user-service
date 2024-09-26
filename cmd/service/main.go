package main

import (
	"fmt"
	user "github.com/s21platform/user-proto/user-proto"
	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/repository/kafka/producer/friends_register"
	"github.com/s21platform/user-service/internal/repository/postgres"
	"github.com/s21platform/user-service/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.MustLoad()

	db := postgres.New(cfg)
	defer db.Close()

	userFriendsRegister := friends_register.New(cfg)

	server := service.New(db, userFriendsRegister)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(infra.UnaryInterceptor),
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
