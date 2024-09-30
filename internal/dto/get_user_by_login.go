package dto

import user "github.com/s21platform/user-proto/user-proto"

type GetUserByLoginIn struct {
	Login string
}

type GetUserByLoginOut struct {
	UUID      string
	IsNewUser bool
}

func ConvertToDTO(in *user.GetUserByLoginIn) *GetUserByLoginIn {
	return &GetUserByLoginIn{
		Login: in.Login,
	}
}

func ConvertFromDTO(in *GetUserByLoginOut) *user.GetUserByLoginOut {
	return &user.GetUserByLoginOut{
		Uuid:      in.UUID,
		IsNewUser: in.IsNewUser,
	}
}
