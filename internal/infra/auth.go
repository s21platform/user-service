package infra

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/s21platform/user-service/internal/config"
)

func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if info.FullMethod == "/UserService/GetUserByLogin" || info.FullMethod == "/UserService/CreateUser" {
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
	return handler(ctx, req)
}

func AuthRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get("X-User-Uuid")
		if requestID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, config.KeyUUID, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
