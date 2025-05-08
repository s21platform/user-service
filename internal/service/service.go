package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/s21platform/friends-proto/friends-proto/new_friend_register"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/model"
	"github.com/s21platform/user-service/pkg/user"
)

type Server struct {
	user.UnimplementedUserServiceServer
	dbRepo     DbRepo
	ufrR       UserFriendsRegisterSrv
	optionhubS OptionhubS
}

func New(repo DbRepo, ufrR UserFriendsRegisterSrv, optionhubService OptionhubS) *Server {
	return &Server{
		dbRepo:     repo,
		ufrR:       ufrR,
		optionhubS: optionhubService,
	}
}

func (s *Server) GetUserByLogin(ctx context.Context, in *user.GetUserByLoginIn) (*user.GetUserByLoginOut, error) {
	m := pkg.FromContext(ctx, config.KeyMetrics)
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	userData, err := s.dbRepo.GetOrSetUserByLogin(in.Login)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user: %v", err))
		return nil, status.Error(codes.NotFound, "Ошибка создания пользователя")
	}
	if userData.IsNew {
		mess := &new_friend_register.NewFriendRegister{
			Email: in.Login,
			Uuid:  userData.Uuid,
		}
		err = s.ufrR.ProduceMessage(ctx, mess, userData.Uuid)
		if err != nil {
			m.Increment("new_friend.error")
			log.Println("error send data to kafka:", err)
		} else {
			m.Increment("new_friend.ok")
		}
	}
	return &user.GetUserByLoginOut{Uuid: userData.Uuid, IsNewUser: userData.IsNew}, nil
}

func (s *Server) IsUserExistByUUID(ctx context.Context, in *user.IsUserExistByUUIDIn) (*user.IsUserExistByUUIDOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	isExist, err := s.dbRepo.IsUserExistByUUID(in.Uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to check user exist: %v", err))
		return nil, status.Errorf(codes.Internal, "Невозможно найти пользователя")
	}
	return &user.IsUserExistByUUIDOut{IsExist: isExist}, nil
}

func (s *Server) GetUsersByUUID(ctx context.Context, in *user.GetUsersByUUIDIn) (*user.GetUsersByUUIDOut, error) {
	_ = ctx
	if in == nil || len(in.UsersUuid) == 0 {
		return nil, fmt.Errorf("no UUIDs provided")
	}

	var userInfoMin []*user.UserInfoMin
	for _, uuid := range in.UsersUuid {
		if uuid.Uuid == "" {
			log.Println("empty UUID provided")
			return &user.GetUsersByUUIDOut{UsersInfo: userInfoMin}, nil
		}

		userInfo, err := s.dbRepo.GetUsersByUUID(uuid.Uuid)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by UUID %s: %w", uuid.Uuid, err)
		}

		userInfoMin = append(userInfoMin, &user.UserInfoMin{
			Uuid:       userInfo.Uuid,
			Login:      userInfo.Login,
			LastAvatar: userInfo.LastAvatar,
			Name:       userInfo.Name,
			Surname:    userInfo.Surname,
		})
	}

	return &user.GetUsersByUUIDOut{UsersInfo: userInfoMin}, nil
}

func (s *Server) GetLoginByUUID(ctx context.Context, in *user.GetLoginByUUIDIn) (*user.GetLoginByUUIDOut, error) {
	login, err := s.dbRepo.GetLoginByUuid(ctx, in.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get login by uuid: %v", err)
	}
	return &user.GetLoginByUUIDOut{Login: login}, nil
}

func (s *Server) GetUserWithOffset(ctx context.Context, in *user.GetUserWithOffsetIn) (*user.GetUserWithOffsetOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, fmt.Errorf("uuid not found in context")
	}
	userWithLimit, total, err := s.dbRepo.GetUserWithLimit(uuid, in.Nickname, in.Limit, in.Offset)
	if err != nil {
		return nil, fmt.Errorf("get user with limit: %w", err)
	}
	var users []*user.User
	for _, u := range userWithLimit {
		users = append(users, &user.User{
			Nickname:   u.Nickname,
			Uuid:       u.UUID,
			AvatarLink: u.AvatarLink,
			Name:       u.Name,
			Surname:    u.Surname,
		})
	}
	return &user.GetUserWithOffsetOut{User: users, Total: total}, nil
}

func (s *Server) GetUserInfoByUUID(ctx context.Context, in *user.GetUserInfoByUUIDIn) (*user.GetUserInfoByUUIDOut, error) {
	_ = ctx.Value(config.KeyUUID).(string)
	// TODO перейти на использование контекстного значения
	userInfo, err := s.dbRepo.GetUserInfoByUUID(ctx, in.Uuid)
	if err != nil {
		log.Println("failed to get user data from repo:", err)
		return nil, status.Errorf(codes.Internal, "failed to get user data from repo")
	}
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", ctx.Value(config.KeyUUID).(string)))
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

func (s *Server) UpdateProfile(ctx context.Context, in *user.UpdateProfileIn) (*user.UpdateProfileOut, error) {
	uuid := ctx.Value(config.KeyUUID).(string)
	var data model.ProfileData
	data.ToDTO(in)
	err := s.dbRepo.UpdateProfile(ctx, data, uuid)
	if err != nil {
		return nil, err
	}
	return &user.UpdateProfileOut{
		Status: true,
	}, nil
}

func (s *Server) SetFriends(ctx context.Context, in *user.SetFriendsIn) (*user.SetFriendsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SetFriends")
	userUUID := ctx.Value(config.KeyUUID).(string)
	if userUUID == "" {
		logger.Error("failed to get user UUID in context")
		return nil, fmt.Errorf("failed to get user UUID in context")
	}
	res, err := s.dbRepo.SetFriends(ctx, userUUID, in.Peer)
	if err != nil {
		logger.Error("failed to SetFriends from BD")
		return nil, err
	}

	if !res {
		logger.Info("user already in friends")
		return &user.SetFriendsOut{Success: false}, nil
	}

	return &user.SetFriendsOut{Success: true}, nil
}

func (s *Server) RemoveFriends(ctx context.Context, in *user.RemoveFriendsIn) (*user.RemoveFriendsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SetFriends")
	userUUID := ctx.Value(config.KeyUUID).(string)
	if userUUID == "" {
		logger.Error("failed to get user UUID in context")
		return nil, fmt.Errorf("failed to get user UUID in context")
	}
	res, err := s.dbRepo.RemoveFriends(ctx, userUUID, in.Peer)
	if err != nil {
		logger.Error("failed to RemoveFriends from BD")
		return nil, err
	}
	if !res {
		logger.Info("user already in friends")
		return &user.RemoveFriendsOut{Success: false}, nil
	}
	return &user.RemoveFriendsOut{Success: true}, nil
}
