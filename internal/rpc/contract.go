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
	GetLoginByUuid(ctx context.Context, uuid string) (string, error)
}

type UserFriendsRegisterSrv interface {
	ProduceMessage(message interface{}) error
}

type OptionhubS interface {
	GetOs(ctx context.Context, id *int64) (*string, error)
}
