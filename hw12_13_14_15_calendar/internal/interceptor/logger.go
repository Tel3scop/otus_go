package interceptor

import (
	"context"
	"time"

	"github.com/Tel3scop/helpers/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// LogInterceptor интерцептор для обработки логов.
func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	interface{},
	error,
) {
	latency := time.Since(time.Now())

	res, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("req", req))
	}

	logger.Info("Request processed",
		zap.String("method", info.FullMethod),
		zap.Any("req", req),
		zap.Any("res", res),
		zap.Duration("duration", latency),
	)

	return res, err
}
