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
)

func RunPprofServer(ctx context.Context, cfg PprofConfig) error {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("listen pprof: %w", err)
	}

	slog.Default().InfoContext(ctx, "pprof server is listening", "address", addr)

	handler := newPprofHandler()

	srv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		srvErr := srv.Serve(lis)
		if srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			errChan <- fmt.Errorf("serve pprof server: %w", srvErr)
			return
		}
	}()

	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		slog.Default().InfoContext(ctx, "shutting down pprof server")

		stopCtx, stop := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer stop()

		err = srv.Shutdown(stopCtx)
		if err != nil {
			return fmt.Errorf("pprof server shutdown: %w", err)
		}

		slog.Default().InfoContext(ctx, "pprof server stopped gracefully")
		return nil
	}
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
