package service

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/s21platform/friends-proto/friends-proto/new_friend_register"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/model"
	"github.com/s21platform/user-service/internal/pkg/generator"
	"github.com/s21platform/user-service/internal/pkg/tx"
	"github.com/s21platform/user-service/pkg/user"
)

type Server struct {
	user.UnimplementedUserServiceServer
	dbRepo     DbRepo
	ufrR       UserFriendsRegisterSrv
	optionhubS OptionhubS
	ucP        UserCreatedProducer
	upcP       UserPostCreatedProduser
}

func New(repo DbRepo, ufrR UserFriendsRegisterSrv, optionhubService OptionhubS, ucP UserCreatedProducer, upcP UserPostCreatedProduser) *Server {
	return &Server{
		dbRepo:     repo,
		ufrR:       ufrR,
		optionhubS: optionhubService,
		ucP:        ucP,
		upcP:       upcP,
	}
}

func (s *Server) GetUserByLogin(ctx context.Context, in *user.GetUserByLoginIn) (*user.GetUserByLoginOut, error) {
	m := pkg.FromContext(ctx, config.KeyMetrics)
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	userData, err := s.dbRepo.GetOrSetUserByLogin(in.Login)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user: %v", err))
		return nil, status.Errorf(codes.NotFound, "Ошибка создания пользователя: %v", err)
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
		return nil, status.Errorf(codes.Internal, "Невозможно найти пользователя: %v", err)
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
		return nil, status.Errorf(codes.Internal, "failed to get user data from repo: %v", err)
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

func (s *Server) CreateUser(ctx context.Context, in *user.CreateUserIn) (*user.CreateUserOut, error) {
	email := strings.TrimSpace(in.Email)
	if email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}
	if len(email) > 100 {
		return nil, status.Errorf(codes.InvalidArgument, "email exceeds maximum length of 100 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`).MatchString(email) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email format")
	}

	userInfo, err := s.dbRepo.GetUserForCreation(ctx, email)
	if err != nil {
		return nil, err
	}

	if userInfo != nil {
		return &user.CreateUserOut{
			UserUuid: userInfo.UUID,
			Nickname: userInfo.Nickname,
		}, nil
	}

	var nickname string
	for {
		nickname = generator.GenerateNickname()
		isAvailable, err := s.dbRepo.CheckNicknameAvailability(ctx, nickname)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to check nickname availability: %v", err)
		}
		if isAvailable {
			break
		}
	}

	userUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate user UUID: %v", err)
	}
	userUUIDStr := userUUID.String()

	err = tx.TxExecute(ctx, func(ctx context.Context) error {
		if err := s.dbRepo.CreateUser(ctx, userUUIDStr, email, nickname); err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}

		if err := s.dbRepo.CreateUserMeta(ctx, userUUIDStr); err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}

		err = s.ucP.ProduceMessage(ctx, user.UserCreatedMessage{
			UserUuid: userUUIDStr,
		}, userUUIDStr)
		if err != nil {
			return fmt.Errorf("failed to produce message: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &user.CreateUserOut{
		UserUuid: userUUIDStr,
		Nickname: nickname,
	}, nil
}

func (s *Server) SetFriends(ctx context.Context, in *user.SetFriendsIn) (*user.SetFriendsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SetFriends")
	userUUID, ok := ctx.Value(config.KeyUUID).(string)

	if !ok || userUUID == "" {
		logger.Error("failed to get user UUID from context")
		return nil, fmt.Errorf("failed to get user UUID from context")
	}

	areFriends, err := s.dbRepo.CheckFriendship(ctx, userUUID, in.Peer)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to check user friendship: %v", err))
		return nil, fmt.Errorf(" failed to check user friendship: %v", err)
	}
	if areFriends {
		logger.Error("user already in friends")
		return &user.SetFriendsOut{Success: true}, nil
	}

	err = s.dbRepo.SetFriends(ctx, userUUID, in.Peer)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to SetFriends from dbRepo: %v", err))
		return nil, err
	}

	return &user.SetFriendsOut{Success: true}, nil
}

func (s *Server) RemoveFriends(ctx context.Context, in *user.RemoveFriendsIn) (*user.RemoveFriendsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("RemoveFriends")
	userUUID, ok := ctx.Value(config.KeyUUID).(string)

	if !ok || userUUID == "" {
		logger.Error("failed to get user UUID from context")
		return nil, fmt.Errorf("failed to get user UUID from context")
	}

	areFriends, err := s.dbRepo.CheckFriendship(ctx, userUUID, in.Peer)
	if err != nil {
		logger.Error(" failed to check user friendship")
		return nil, fmt.Errorf(" failed to check user friendshipt: %v", err)
	}
	if !areFriends {
		logger.Error("user already not friends")
		return &user.RemoveFriendsOut{Success: false}, nil
	}

	err = s.dbRepo.RemoveFriends(ctx, userUUID, in.Peer)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to RemoveFriends from dbRepo: %v", err))
		return nil, err
	}
	return &user.RemoveFriendsOut{Success: true}, nil
}

func (s *Server) GetCountFriends(ctx context.Context, in *user.EmptyFriends) (*user.GetCountFriendsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetCountFriends")
	userUUID, ok := ctx.Value(config.KeyUUID).(string)

	if !ok || userUUID == "" {
		logger.Error("failed to get user UUID from context")
		return nil, fmt.Errorf("failed to get user UUID from context")
	}

	subscription, err := s.dbRepo.GetSubscriptionCount(ctx, userUUID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get subscription count: %v", err))
		return nil, err
	}
	subscribers, err := s.dbRepo.GetSubscribersCount(ctx, userUUID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get subscribers count: %v", err))
		return nil, err
	}
	return &user.GetCountFriendsOut{Subscription: subscription, Subscribers: subscribers}, nil
}

func (s *Server) GetPeerFollow(ctx context.Context, in *user.GetPeerFollowIn) (*user.GetPeerFollowOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetPeerFollow")

	userUUID, ok := ctx.Value(config.KeyUUID).(string)
	if !ok || userUUID == "" {
		logger.Error("failed to get user UUID from context")
		return nil, fmt.Errorf("failed to get user UUID from context")
	}

	follow, err := s.dbRepo.GetPeerFollow(ctx, in.Uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get peer follow: %v", err))
		return nil, err
	}
	peers := make([]*user.Peer, 0)

	for _, peerUuid := range follow {
		peers = append(peers, &user.Peer{Uuid: peerUuid})
	}

	return &user.GetPeerFollowOut{Subscription: peers}, nil
}

func (s *Server) GetWhoFollowPeer(ctx context.Context, in *user.GetWhoFollowPeerIn) (*user.GetWhoFollowPeerOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetWhoFollowPeer")
	userUUID, ok := ctx.Value(config.KeyUUID).(string)

	if !ok || userUUID == "" {
		logger.Error("failed to get user UUID from context")
		return nil, fmt.Errorf("failed to get user UUID from context")
	}

	follow, err := s.dbRepo.GetWhoFollowPeer(ctx, in.Uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get peer follow: %v", err))
		return nil, err
	}
	peers := make([]*user.Peer, 0)
	for _, uuid := range follow {
		peers = append(peers, &user.Peer{Uuid: uuid})
	}
	return &user.GetWhoFollowPeerOut{Subscribers: peers}, nil
}

func (s *Server) CheckFriendship(ctx context.Context, in *user.CheckFriendshipIn) (*user.CheckFriendshipOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("CheckFriendship")
	userUUID, ok := ctx.Value(config.KeyUUID).(string)

	if !ok || userUUID == "" {
		logger.Error("failed to get user UUID from context")
		return nil, fmt.Errorf("failed to get user UUID from context")
	}

	succses, err := s.dbRepo.CheckFriendship(ctx, userUUID, in.Uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to check user friendship: %v", err))
		return nil, err
	}
	return &user.CheckFriendshipOut{
		Succses: succses,
	}, nil
}

func (s *Server) CreatePost(ctx context.Context, in *user.CreatePostIn) (*user.CreatePostOut, error) {
	ownerUUIDStr, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to retrieve uuid")
	}

	ownerUUID, err := uuid.Parse(ownerUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid: %v", err)
	}

	newPostUUID, err := s.dbRepo.CreatePost(ctx, ownerUUID, in.Content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %v", err)
	}

	msg := &user.UserPostCreated{
		UserUuid: ownerUUIDStr,
		PostId:   newPostUUID,
	}

	err = s.upcP.ProduceMessage(ctx, msg, ownerUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to produce message: %v", err)
	}

	return &user.CreatePostOut{PostUuid: newPostUUID}, nil
}

func (s *Server) GetPostsByIds(ctx context.Context, in *user.GetPostsByIdsIn) (*user.GetPostsByIdsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetPostsByIds")

	userUUID, ok := ctx.Value(config.KeyUUID).(string)
	if !ok || userUUID == "" {
		logger.Error("user UUID required")
		return nil, fmt.Errorf("failed to get user UUID from context")
	}

	posts, err := s.dbRepo.GetPostsByIds(ctx, in.PostUuids)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get posts by ids: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get posts from db: %v", err)
	}
	return &user.GetPostsByIdsOut{Posts: posts.FromDTO()}, nil
}
