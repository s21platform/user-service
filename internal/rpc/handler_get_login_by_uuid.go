package rpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	user_proto "github.com/s21platform/user-proto/user-proto"
)

func (s *Server) GetLoginByUUID(ctx context.Context, in *user_proto.GetLoginByUUIDIn) (*user_proto.GetLoginByUUIDOut, error) {
	login, err := s.dbRepo.GetLoginByUuid(ctx, in.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get login by uuid: %v", err)
	}
	return &user_proto.GetLoginByUUIDOut{Login: login}, nil
}
