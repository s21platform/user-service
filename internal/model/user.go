package model

import (
	"time"

	"github.com/samber/lo"

	optionhub_lib "github.com/s21platform/optionhub-lib"

	"github.com/s21platform/user-service/pkg/user"
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
	Name      string     `db:"name"`
	Birthdate *time.Time `db:"birthdate"`
	Telegram  string     `db:"telegram"`
	Git       string     `db:"git"`
	OsId      int64      `db:"os_id"`
}

func (pd *ProfileData) ToDTOv2(parsedData []optionhub_lib.AttributeValue) {
	for _, v := range parsedData {
		if f, ok := AttributeSetters[v.AttributeId]; ok {
			f(v, pd)
		}
	}
}

func (pd *ProfileData) ToDTO(in *user.UpdateProfileIn) {
	birthdate, err := time.Parse(time.RFC3339, in.Birthday)
	bd := lo.Ternary(err != nil, nil, &birthdate)
	pd.Name = in.Name
	pd.Birthdate = bd
	pd.Telegram = in.Telegram
	pd.Git = in.Github
	pd.OsId = in.OsId
}

type UserWithLimit struct {
	Nickname   string `db:"login"`
	UUID       string `db:"uuid"`
	AvatarLink string `db:"last_avatar_link"`
	Name       string `db:"name"`
	Surname    string `db:"surname"`
}
