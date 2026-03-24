package internal

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func newGRPCServerOptions(cfg GRPCConfig) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.MaxSendMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.MaxConnectionIdle,
			MaxConnectionAge:  cfg.MaxConnectionAge,
		}),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(loggingHandler()),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(),
			logging.StreamServerInterceptor(loggingHandler()),
		),
	}
}

func loggingHandler() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelError:
			slog.ErrorContext(ctx, msg, fields...)
		case logging.LevelWarn:
			slog.WarnContext(ctx, msg, fields...)
		case logging.LevelInfo:
			slog.InfoContext(ctx, msg, fields...)
		case logging.LevelDebug:
			slog.DebugContext(ctx, msg, fields...)
		}
	})
}
