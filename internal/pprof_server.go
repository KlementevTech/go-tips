package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"

	"golang.org/x/sync/errgroup"
)

func RunPprofServer(ctx context.Context, g *errgroup.Group, cfg PprofConfig) error {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return err
	}

	slog.Default().InfoContext(ctx, "pprof server is listening", "address", addr)

	srv := &http.Server{
		Handler:           newPprofHandler(),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	g.Go(func() error {
		serveErr := srv.Serve(lis)
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			return serveErr
		}

		slog.Default().InfoContext(ctx, "pprof server is stopped")
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		stopCtx, stop := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer stop()

		if stopErr := srv.Shutdown(stopCtx); stopErr != nil {
			return fmt.Errorf("shutdown pprof server: %w", stopErr)
		}
		return nil
	})

	return nil
}

func newPprofHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return handler
}
