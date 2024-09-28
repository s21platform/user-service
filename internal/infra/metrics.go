package infra

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/s21platform/metrics-lib/pkg"
)

func MetricsInterceptor(metrics *pkg.Metrics) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		t := time.Now()
		method := strings.Trim(strings.ReplaceAll(info.FullMethod, "/", "_"), "_")
		metrics.Increment(method)
		resp, err := handler(ctx, req)
		if err != nil {
			metrics.Increment(method + "_error")
		}
		metrics.Duration(time.Since(t).Milliseconds(), method)
		return resp, err
	}
}
