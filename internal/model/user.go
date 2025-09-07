package model

import (
	"time"
)

type UserInfo struct {
	Nickname       string     `db:"login"`
	LastAvatarLink string     `db:"last_avatar_link"`
	Name           *string    `db:"name"`
	Surname        *string    `db:"surname"`
	Birthdate      *time.Time `db:"birthdate"`
	Phone          *string    `db:"phone"`
	Telegram       *string    `db:"telegram"`
	Git            *string    `db:"git"`
	CityId         *int64     `db:"city_id"`
	OSId           *int64     `db:"os_id"`
	WorkId         *int64     `db:"work_id"`
	UniversityId   *int64     `db:"university_id"`
}

type UserInfoMin struct {
	Uuid       string `db:"uuid"`
	Login      string `db:"login"`
	LastAvatar string `db:"last_avatar_link"`
	Name       string `db:"name"`
	Surname    string `db:"surname"`
}

type ProfileData struct {
	Name     *string    `db:"name"`
	Surname  *string    `db:"surname"`
	Birthday *time.Time `db:"birthdate"`
	CityId   *int64     `db:"city_id"`
	Telegram *string    `db:"telegram"`
}

type Personality struct {
	Name      *string    `db:"name"`
	Surname   *string    `db:"surname"`
	Birthdate *time.Time `db:"birthdate"`
}

type UserAttributes struct {
	Name      *string    `db:"name"`
	Surname   *string    `db:"surname"`
	Birthdate *time.Time `db:"birthdate"`
	CityId    *int64     `db:"city_id"`
	Telegram  *string    `db:"telegram"`
}

type UserWithLimit struct {
	Nickname   string `db:"login"`
	UUID       string `db:"uuid"`
	AvatarLink string `db:"last_avatar_link"`
	Name       string `db:"name"`
	Surname    string `db:"surname"`
}

//type Peer struct {
//	uuid string
//}
