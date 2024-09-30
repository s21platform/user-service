package infra

import (
	"context"
	"log"
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
		log.Printf("second interceptor for: %s", info.FullMethod)
		t := time.Now()
		method := strings.Trim(strings.ReplaceAll(info.FullMethod, "/", "_"), "_")
		metrics.Increment(method)
		log.Printf("continue second interceptor before call: %s", info.FullMethod)
		resp, err := handler(ctx, req)
		log.Printf("continue second interceptor after call: %s", info.FullMethod)

		if err != nil {
			metrics.Increment(method + "_error")
		}
		metrics.Duration(time.Since(t).Milliseconds(), method)
		log.Printf("end second interceptor: %s", info.FullMethod)
		return resp, err
	}
}
