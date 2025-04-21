package rpc

import (
	"context"
	"fmt"
	"strconv"

	user "github.com/s21platform/user-proto/user-proto"
	"github.com/s21platform/user-service/internal/config"
)

func (s *Server) GetUsersInfoWithOffset(ctx context.Context, in *user.GetUserWithOffsetIn) (*user.GetUserWithOffsetOutAll, error) {
	fmt.Println("get info offset", in.Offset)
	_, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get uuid from context")
	}
	res, total, err := s.dbRepo.GetAllInfoUsers(ctx, in.Nickname, in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}
	users := []*user.GetUserInfoByUUIDOut{}
	for _, u := range res {
		var bDay *string
		if u.Birthdate != nil {
			formatted := u.Birthdate.Format("2006-01-02")
			bDay = &formatted
		}
		users = append(users, &user.GetUserInfoByUUIDOut{
			Nickname:  u.Nickname,
			Avatar:    u.LastAvatarLink,
			Name:      u.Name,
			Surname:   u.Surname,
			Birthdate: bDay,
			Phone:     u.Phone,
			Telegram:  u.Telegram,
			Git:       u.Git,
			City:      int64PtrToStringPtr(u.CityId),
			Os: &user.GetOs{
				Label: func() string {
					if u.OSId == nil {
						return ""
					}
					return strconv.FormatInt(*u.OSId, 10)
				}(),
			},
			Work:       int64PtrToStringPtr(u.WorkId),
			University: int64PtrToStringPtr(u.UniversityId),
			Uuid:       u.UUID,
		})
	}
	return &user.GetUserWithOffsetOutAll{User: users, Total: total}, nil
}

func int64PtrToStringPtr(id *int64) *string {
	if id == nil {
		return nil
	}
	s := strconv.FormatInt(*id, 10)
	return &s
}
