package rpc

import (
	"context"
	"fmt"
	user "github.com/s21platform/user-proto/user-proto"
	"github.com/s21platform/user-service/internal/config"
)

func (s *Server) GetUserWithLimit(ctx context.Context, in *user.GetUserWithLimitIn) (*user.GetUserWithLimitOut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, fmt.Errorf("uuid not found in context")
	}
	userWithLimit, err := s.dbRepo.GetUserWithLimit(uuid, in.Limit, in.Offset)
	if err != nil {
		return nil, fmt.Errorf("get user with limit: %w", err)
	}
	var users []*user.User
	for _, u := range userWithLimit {
		users = append(users, &user.User{
			Nickname:   u.Nickname,
			Uuid:       u.UUID,
			AvatarLink: u.Avatar_link,
		})
	}
	return &user.GetUserWithLimitOut{User: users, Total: int64(len(users))}, nil
}
