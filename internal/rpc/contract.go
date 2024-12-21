package rpc

import (
	"context"

	"github.com/s21platform/user-service/internal/model"
	"github.com/s21platform/user-service/internal/repository/postgres"
)

type DbRepo interface {
	IsUserExistByUUID(uuid string) (bool, error)
	GetOrSetUserByLogin(email string) (*postgres.CheckUser, error)
	GetUserInfoByUUID(ctx context.Context, uuid string) (model.UserInfo, error)
	GetUserWithLimit(uuid, nicName string, limit int64, offset int64) ([]model.UserWithLimit, error)
	GetLoginByUuid(ctx context.Context, uuid string) (string, error)
	UpdateProfile(ctx context.Context, data model.ProfileData, userUuid string) error
}

type UserFriendsRegisterSrv interface {
	ProduceMessage(message interface{}) error
}

type OptionhubS interface {
	GetOs(ctx context.Context, id *int64) (*model.OS, error)
}
