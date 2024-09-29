package rpc

import (
	"context"
	"log"

	user "github.com/s21platform/user-proto/user-proto"
	"github.com/s21platform/user-service/internal/config"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUserInfoByUUID(ctx context.Context, in *user.GetUserInfoByUUIDIn) (*user.GetUserInfoByUUIDOut, error) {
	test := ctx.Value(config.KeyUUID)
	t, ok := test.(string)
	if !ok {
		log.Println("GetUserInfoByUUID error:", t)
	}
	log.Println("uuid from context:", t)
	userInfo, err := s.dbRepo.GetUserInfoByUUID(ctx, in.Uuid)
	if err != nil {
		log.Println("failed to get user data from repo:", err)
		return nil, status.Errorf(codes.Internal, "failed to get user data from repo")
	}
	resp := &user.GetUserInfoByUUIDOut{
		Nickname:   userInfo.Nickname,
		Avatar:     userInfo.LastAvatarLink,
		Name:       userInfo.Name,
		Surname:    userInfo.Surname,
		Birthdate:  userInfo.Birthdate,
		Phone:      userInfo.Phone,
		Telegram:   userInfo.Telegram,
		Git:        userInfo.Git,
		City:       lo.ToPtr("Москва [HC]"),
		Os:         lo.ToPtr("Mac OS [HC]"),
		Work:       lo.ToPtr("Avito tech [HC]"),
		University: lo.ToPtr("НИУ МЭИ [HC]"),
	}
	return resp, nil
}
