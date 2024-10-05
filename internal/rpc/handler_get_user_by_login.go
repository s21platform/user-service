package rpc

import (
	"context"
	"log"

	"github.com/s21platform/user-service/internal/config"

	"github.com/s21platform/metrics-lib/pkg"

	user "github.com/s21platform/user-proto/user-proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUserByLogin(ctx context.Context, in *user.GetUserByLoginIn) (*user.GetUserByLoginOut, error) {
	m := pkg.FromContext(ctx, config.KeyMetrics)
	userData, err := s.dbRepo.GetOrSetUserByLogin(in.Login)
	if err != nil {
		log.Println("GetUserByLogin error:", err)
		return nil, status.Error(codes.NotFound, "Ошибка создания пользователя")
	}
	if userData.IsNew {
		err = s.ufrR.SendMessage(ctx, in.Login, userData.Uuid)
		if err != nil {
			m.Increment("new_friend.error")
			log.Println("error send data to kafka:", err)
		} else {
			m.Increment("new_friend.ok")
		}
	}

	return &user.GetUserByLoginOut{Uuid: userData.Uuid, IsNewUser: userData.IsNew}, nil
}
