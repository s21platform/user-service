package rpc

import (
	"context"
	"log"
	"time"

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

	osInfo, err := s.optionhubS.GetOs(ctx, userInfo.OSId)
	if err != nil {
		log.Printf("cannot get os, err: %v\n", err)
	}

	var birthday *string
	if userInfo.Birthdate != nil {
		birthday = lo.ToPtr(userInfo.Birthdate.Format(time.DateOnly))
	}

	var os *user.GetOs
	if osInfo != nil {
		os = &user.GetOs{
			Id:    osInfo.Id,
			Label: osInfo.Label,
		}
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
