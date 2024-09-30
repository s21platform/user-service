package rpc

import (
	"context"
	"log"

	user "github.com/s21platform/user-proto/user-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUserByLogin(ctx context.Context, in *user.GetUserByLoginIn) (*user.GetUserByLoginOut, error) {
	userData, err := s.dbRepo.GetOrSetUserByLogin(in.Login)
	if err != nil {
		log.Println("GetUserByLogin error:", err)
		return nil, status.Error(codes.NotFound, "Ошибка создания пользователя")
	}
	if userData.IsNew {
		err = s.ufrR.SendMessage(ctx, in.Login, userData.Uuid)
		if err != nil {
			log.Println("error send data to kafka:", err)
			// FIXME Тут не надо возвращать ошибку! она заблочит нормальную работу в случае неполадок
			//return nil, status.Error(codes.Unknown, "Ошибка отправки в очередь")
		}
	}
	return &user.GetUserByLoginOut{Uuid: userData.Uuid, IsNewUser: userData.IsNew}, nil
}
