package model

type UserAuthInfo struct {
	UUID     string `db:"uuid"`
	Nickname string `db:"login"`
}
