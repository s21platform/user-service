package service

import "github.com/s21platform/user-service/internal/repository/postgres"

type DbRepo interface {
	IsUserExistByUUID(uuid string) (bool, error)
	GetOrSetUserByLogin(email string) (*postgres.CheckUser, error)
}
