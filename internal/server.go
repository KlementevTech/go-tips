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
	lis, err := lc.Listen(ctx, "tcp4", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	slog.Default().InfoContext(ctx, "pprof server is listening", "address", addr)

	g.Go(func() error {
		handler := pprofHandler()

		srv := &http.Server{
			Handler:           handler,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		}

		errChan := make(chan error, 1)
		defer close(errChan)

		go func() {
			serveErr := srv.Serve(lis)
			if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
				errChan <- serveErr
			}
		}()

		select {
		case err = <-errChan:
			return fmt.Errorf("failed to serve: %w", err)
		case <-ctx.Done():
			stopCtx, stop := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
			defer stop()

			if stopErr := srv.Shutdown(stopCtx); stopErr != nil {
				return fmt.Errorf("could not shutdown pprof server: %w", stopErr)
			}

			slog.Default().InfoContext(ctx, "pprof server is stopped")
			return nil
		}
	})

	return nil
}

func pprofHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return handler
}
