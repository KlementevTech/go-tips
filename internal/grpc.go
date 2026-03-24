package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func RunGRPCServer(ctx context.Context, cfg GRPCConfig) error {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("listen gRPC: %w", err)
	}

	slog.Default().InfoContext(ctx, "gRPC server is listening", "address", addr)

	srv := grpc.NewServer(newGRPCServerOptions(cfg)...)

	if cfg.EnableHealth {
		registerHealthServer(srv)
		slog.Default().InfoContext(ctx, "gRPC health checking enabled")
	}

	if cfg.EnableReflection {
		reflection.Register(srv)
		slog.Default().InfoContext(ctx, "gRPC server reflection enabled")
	}

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		srvErr := srv.Serve(lis)
		if srvErr != nil && !errors.Is(srvErr, grpc.ErrServerStopped) {
			errChan <- fmt.Errorf("serve gRPC server: %w", srvErr)
			return
		}
	}()

	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		slog.Default().InfoContext(ctx, "shutting down gRPC server")

		stopCtx, stop := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer stop()

		err = stopGRPC(stopCtx, srv)
		if err != nil {
			return err
		}

		slog.Default().InfoContext(ctx, "gRPC server stopped gracefully")
		return nil
	}
}

func stopGRPC(ctx context.Context, grpcServer *grpc.Server) error {
	done := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		grpcServer.Stop()
		return errors.New("gRPC shutdown timeout, forcing stop")
	}
}

func registerHealthServer(srv *grpc.Server) {
	hs := health.NewServer()
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(srv, hs)
}
