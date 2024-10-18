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
	_ = ctx.Value(config.KeyUUID).(string)
	// TODO перейти на использование контекстного значения
	userInfo, err := s.dbRepo.GetUserInfoByUUID(ctx, in.Uuid)
	if err != nil {
		log.Println("failed to get user data from repo:", err)
		return nil, status.Errorf(codes.Internal, "failed to get user data from repo")
	}

	var birthday *user.Birthday
	if userInfo.Birthdate != nil {
		birthday = &user.Birthday{
			Day:   int64(userInfo.Birthdate.Day()),
			Month: int64(userInfo.Birthdate.Month()),
			Year:  int64(userInfo.Birthdate.Year()),
		}
	}

	os, err := s.optionhubS.GetOs(ctx, userInfo.OSId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get os name")
	}

	resp := &user.GetUserInfoByUUIDOut{
		Nickname:   userInfo.Nickname,
		Avatar:     userInfo.LastAvatarLink,
		Name:       userInfo.Name,
		Surname:    userInfo.Surname,
		Birthdate:  birthday,
		Phone:      userInfo.Phone,
		Telegram:   userInfo.Telegram,
		Git:        userInfo.Git,
		City:       lo.ToPtr("Москва [HC]"),
		Os:         os,
		Work:       lo.ToPtr("Avito tech [HC]"),
		University: lo.ToPtr("НИУ МЭИ [HC]"),
	}
	return resp, nil
}
