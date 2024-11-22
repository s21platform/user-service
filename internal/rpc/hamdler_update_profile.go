package rpc

import (
	"context"

	user_proto "github.com/s21platform/user-proto/user-proto"
	"github.com/s21platform/user-service/internal/config"
	"github.com/s21platform/user-service/internal/model"
)

func (s *Server) UpdateProfile(ctx context.Context, in *user_proto.UpdateProfileIn) (*user_proto.UpdateProfileOut, error) {
	uuid := ctx.Value(config.KeyUUID).(string)
	var data model.ProfileData
	data.ToDTO(in)
	err := s.dbRepo.UpdateProfile(ctx, data, uuid)
	if err != nil {
		return nil, err
	}
	return &user_proto.UpdateProfileOut{
		Status: true,
	}, nil
}
