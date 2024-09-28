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
}

type UserFriendsRegisterSrv interface {
	SendMessage(ctx context.Context, email string, uuid string) error
}
