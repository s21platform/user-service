package rpc

import (
	"context"
	"fmt"
	"log"

	user "github.com/s21platform/user-proto/user-proto"
)

func (s *Server) GetUsersByUUID(ctx context.Context, in *user.GetUsersByUUIDIn) (*user.GetUsersByUUIDOut, error) {
	_ = ctx
	if in == nil || len(in.UsersUuid) == 0 {
		return nil, fmt.Errorf("no UUIDs provided")
	}

	var userInfoMin []*user.UserInfoMin
	for _, uuid := range in.UsersUuid {
		if uuid.Uuid == "" {
			log.Println("empty UUID provided")
			return &user.GetUsersByUUIDOut{UsersInfo: userInfoMin}, nil
		}

		userInfo, err := s.dbRepo.GetUsersByUUID(uuid.Uuid)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by UUID %s: %w", uuid.Uuid, err)
		}

		userInfoMin = append(userInfoMin, &user.UserInfoMin{
			Uuid:       userInfo.Uuid,
			Login:      userInfo.Login,
			LastAvatar: userInfo.LastAvatar,
			Name:       userInfo.Name,
			Surname:    userInfo.Surname,
		})
	}

	return &user.GetUsersByUUIDOut{UsersInfo: userInfoMin}, nil
}
