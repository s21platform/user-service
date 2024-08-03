package service

import (
	"context"
	user_proto "github.com/s21platform/user-proto/user-proto"
	"github.com/s21platform/user-service/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	user_proto.UnimplementedUserServiceServer
	dbRepo DbRepo
}

func New(cfg *config.Config, repo DbRepo) *Server {
	return &Server{dbRepo: repo}
}

func (s *Server) GetUserByLogin(ctx context.Context, in *user_proto.GetUserByLoginIn) (*user_proto.GetUserByLoginOut, error) {
	user, err := s.dbRepo.GetOrSetUserByLogin(in.Login)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Ошибка создания пользователя")
	}
	return &user_proto.GetUserByLoginOut{Uuid: user.Uuid, IsNewUser: user.IsNew}, nil
}

func (s *Server) IsUserExistByUUID(ctx context.Context, in *user_proto.IsUserExistByUUIDIn) (*user_proto.IsUserExistByUUIDOut, error) {
	isExist, err := s.dbRepo.IsUserExistByUUID(in.Uuid)
	if err != nil {
		log.Println("IsUserExistByUUID error:", err)
		return nil, status.Errorf(codes.Internal, "Невозможно найти пользователя")
	}
	return &user_proto.IsUserExistByUUIDOut{IsExist: isExist}, nil
}
