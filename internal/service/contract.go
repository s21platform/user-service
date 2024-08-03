package service

type DbRepo interface {
	IsUserExistByUUID(uuid string) (bool, error)
}
