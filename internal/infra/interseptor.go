package infra

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if info.FullMethod == "/UserService/GetUserByLogin" {
		return handler(ctx, req)
	}
	if ok {
		if userIDs := md["uuid"]; len(userIDs) > 0 {
			ctx = context.WithValue(ctx, "uuid", userIDs[0])
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "no uuid found in metadata")
		}
	} else {
		return nil, status.Errorf(codes.Unauthenticated, "no uuid found in metadata")
	}

	return handler(ctx, req)
}
