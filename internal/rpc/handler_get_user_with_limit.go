package rpc

import (
	"context"
	"fmt"

	user "github.com/s21platform/user-proto/user-proto"
	"github.com/s21platform/user-service/internal/config"
)

func (s *Server) GetUserWithOffset(ctx context.Context, in *user.GetUserWithOffsetIn) (*user.GetUserWithOffsetOut, error) {
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
			Name:       u.Name,
			Surname:    u.Surname,
		})
	}
	return &user.GetUserWithOffsetOut{User: users, Total: int64(len(users))}, nil
}
