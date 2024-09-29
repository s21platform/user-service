package infra

import (
	"context"
	"log"

	"github.com/s21platform/user-service/internal/config"

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
	log.Printf("first interceptor: %s", info.FullMethod)
	if info.FullMethod == "/UserService/GetUserByLogin" {
		log.Printf("skip first interceptor for: %s", info.FullMethod)
		return handler(ctx, req)
	}
	if ok {
		if userIDs := md["uuid"]; len(userIDs) > 0 {
			ctx = context.WithValue(ctx, config.KeyUUID, userIDs[0])
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "no uuid found in metadata")
		}
	} else {
		return nil, status.Errorf(codes.Unauthenticated, "no uuid found in metadata")
	}
	log.Printf("continue first interceptor for: %s", info.FullMethod)
	return handler(ctx, req)
}
