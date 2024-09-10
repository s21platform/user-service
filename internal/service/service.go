package service

import (
	"context"
	user "github.com/s21platform/user-proto/user-proto"
	"github.com/s21platform/user-service/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	user.UnimplementedUserServiceServer
	dbRepo DbRepo
	ufrR   UserFriendsRegisterSrv
}

func New(cfg *config.Config, repo DbRepo, ufrR UserFriendsRegisterSrv) *Server {
	return &Server{dbRepo: repo, ufrR: ufrR}
}

func (s *Server) GetUserByLogin(ctx context.Context, in *user.GetUserByLoginIn) (*user.GetUserByLoginOut, error) {
	userData, err := s.dbRepo.GetOrSetUserByLogin(in.Login)
	if err != nil {
		log.Println("GetUserByLogin error:", err)
		return nil, status.Error(codes.NotFound, "Ошибка создания пользователя")
	}
	if userData.IsNew {
		err = s.ufrR.SendMessage(ctx, in.Login)
		if err != nil {
			log.Println("error send data to kafka:", err)
			// FIXME Тут не надо возвращать ошибку! она заблочит нормальную работу в случае неполадок
			return nil, status.Error(codes.Unknown, "Ошибка отправки в очередь")
		}
	}
	return &user.GetUserByLoginOut{Uuid: userData.Uuid, IsNewUser: userData.IsNew}, nil
}

func (s *Server) IsUserExistByUUID(ctx context.Context, in *user.IsUserExistByUUIDIn) (*user.IsUserExistByUUIDOut, error) {
	isExist, err := s.dbRepo.IsUserExistByUUID(in.Uuid)
	if err != nil {
		log.Println("IsUserExistByUUID error:", err)
		return nil, status.Errorf(codes.Internal, "Невозможно найти пользователя")
	}
	return &user.IsUserExistByUUIDOut{IsExist: isExist}, nil
}
