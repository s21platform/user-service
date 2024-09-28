package model

type UserInfo struct {
	Nickname       string  `db:"login"`
	LastAvatarLink string  `db:"last_avatar_link"`
	Name           *string `db:"name"`
	Surname        *string `db:"surname"`
	Birthdate      *string `db:"birthdate"`
	Phone          *string `db:"phone"`
	Telegram       *string `db:"telegram"`
	Git            *string `db:"git"`
	CityId         *string `db:"city_id"`
	OSId           *string `db:"os_id"`
	WorkId         *string `db:"work_id"`
	UniversityId   *string `db:"university_id"`
}
