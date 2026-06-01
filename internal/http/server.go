package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func RunServer(ctx context.Context, addr string, handler http.Handler, opt ...ServerOption) error {
	opts := defaultServerOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	logger := opts.logger

	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("listen http: %w", err)
	}

	httpSrv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: opts.readHeaderTimeout,
	}

	errChan := make(chan error, 1)
	go func() {
		logger.InfoContext(ctx, "starting http server", "address", addr)

		serveErr := httpSrv.Serve(lis)
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errChan <- serveErr
		}
	}()

	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		err = shutdown(httpSrv, opts.shutdownTimeout)
		if err != nil {
			logger.ErrorContext(ctx, "shutting down http server", slog.String("error", err.Error()))
		} else {
			logger.InfoContext(ctx, "http server stopped gracefully")
		}

		return nil
	}
}

func shutdown(server *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return server.Shutdown(ctx)
}
