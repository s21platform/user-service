package rpc

import (
	"context"
	"log"

	user "github.com/s21platform/user-proto/user-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) IsUserExistByUUID(ctx context.Context, in *user.IsUserExistByUUIDIn) (*user.IsUserExistByUUIDOut, error) {
	isExist, err := s.dbRepo.IsUserExistByUUID(in.Uuid)
	if err != nil {
		log.Println("IsUserExistByUUID error:", err)
		return nil, status.Errorf(codes.Internal, "Невозможно найти пользователя")
	}
	return &user.IsUserExistByUUIDOut{IsExist: isExist}, nil
}
