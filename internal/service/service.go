package service

import (
	"context"
	user "github.com/s21platform/user-proto/user-proto"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	user.UnimplementedUserServiceServer
	dbRepo DbRepo
	ufrR   UserFriendsRegisterSrv
}

func New(repo DbRepo, ufrR UserFriendsRegisterSrv) *Server {
	return &Server{dbRepo: repo, ufrR: ufrR}
}

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

func (s *Server) GetUserInfoByUUID(ctx context.Context, in *user.GetUserInfoByUUIDIn) (*user.GetUserInfoByUUIDOut, error) {
	test := ctx.Value("uuid")
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
