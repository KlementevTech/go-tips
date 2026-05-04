package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/KlementevTech/gotips/internal/api/middleware/errlog"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Address           string        `mapstructure:"address"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	MaxRecvMsgSize    int           `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSize    int           `mapstructure:"max_send_msg_size"`
	MaxConnectionIdle time.Duration `mapstructure:"max_connection_idle"`
	MaxConnectionAge  time.Duration `mapstructure:"max_connection_age"`
	EnableReflection  bool          `mapstructure:"enable_reflection"`
	EnableHealth      bool          `mapstructure:"enable_health"`
}

func RunServer(ctx context.Context, cfg Config, register func(s *grpc.Server) error) error {
	middlewares := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			errlog.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(),
			errlog.StreamServerInterceptor(),
		),
	}

	options := append([]grpc.ServerOption{
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.MaxSendMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.MaxConnectionIdle,
			MaxConnectionAge:  cfg.MaxConnectionAge,
		}),
	}, middlewares...)

	srv := grpc.NewServer(options...)

	if cfg.EnableHealth {
		registerHealthServer(srv)
	}

	if cfg.EnableReflection {
		reflection.Register(srv)
	}

	err := register(srv)
	if err != nil {
		return fmt.Errorf("register gRPC servers: %w", err)
	}

	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("listen gRPC: %w", err)
	}

	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		slog.Default().InfoContext(ctx, "starting gRPC server",
			slog.String("address", cfg.Address),
			slog.Bool("reflection", cfg.EnableReflection),
			slog.Bool("health_check", cfg.EnableHealth),
		)

		err = srv.Serve(lis)
		if err != nil {
			errCh <- fmt.Errorf("serve gRPC server: %w", err)
		}
	}()

	select {
	case err = <-errCh:
		return err
	case <-ctx.Done():
		select {
		case err = <-errCh:
			return err
		default:
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		shutdownGrpc(shutdownCtx, srv)
		slog.Default().InfoContext(ctx, "gRPC server stopped gracefully")
	}
	return nil
}

func shutdownGrpc(ctx context.Context, grpcServer *grpc.Server) {
	done := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		grpcServer.Stop()
		slog.Default().ErrorContext(ctx, "gRPC shutdown timeout, forcing stop")
	}
}

func registerHealthServer(srv *grpc.Server) {
	hs := health.NewServer()
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(srv, hs)
}
