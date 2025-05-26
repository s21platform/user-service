//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	"github.com/s21platform/user-service/internal/model"
	"github.com/s21platform/user-service/internal/repository/postgres"
)

type DbRepo interface {
	IsUserExistByUUID(uuid string) (bool, error)
	GetOrSetUserByLogin(email string) (*postgres.CheckUser, error)
	GetUserInfoByUUID(ctx context.Context, uuid string) (model.UserInfo, error)
	GetUserWithLimit(uuid, nickname string, limit int64, offset int64) ([]model.UserWithLimit, int64, error)
	GetLoginByUuid(ctx context.Context, uuid string) (string, error)
	UpdateProfile(ctx context.Context, data model.ProfileData, userUuid string) error
	GetUsersByUUID(uuid string) (model.UserInfoMin, error)
	SetFriends(ctx context.Context, peer1, peer2 string) error
	RemoveFriends(ctx context.Context, peer1, peer2 string) error
	CheckFriendship(ctx context.Context, peer1, peer2 string) (bool, error)
	GetSubscribersCount(ctx context.Context, userUUID string) (int64, error)
	GetSubscriptionCount(ctx context.Context, userUUID string) (int64, error)
	GetPeerFollow(ctx context.Context, userUUID string) ([]string, error)
	GetWhoFollowPeer(ctx context.Context, userUUID string) ([]string, error)
	CheckNicknameAvailability(ctx context.Context, nickname string) (bool, error)
	CreateUser(ctx context.Context, userUUID string, email string, nickname string) error
	GetUserForCreation(ctx context.Context, email string) (*model.UserAuthInfo, error)
	CreatePost(ctx context.Context, uuid, content string) (string, error)
}

type UserFriendsRegisterSrv interface {
	ProduceMessage(ctx context.Context, message any, key any) error
}

type UserCreatedProducer interface {
	ProduceMessage(ctx context.Context, message any, key any) error
}

type OptionhubS interface {
	GetOs(ctx context.Context, id *int64) (*model.OS, error)
}
