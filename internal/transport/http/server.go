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

const (
	shutdownTimeout   = time.Second * 30
	readHeaderTimeout = time.Second * 10
)

type Config struct {
	Alias             string
	Addr              string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

func NewDefaultConfig(addr, alias string) *Config {
	return &Config{
		Alias:             alias,
		Addr:              addr,
		ReadHeaderTimeout: readHeaderTimeout,
		ShutdownTimeout:   shutdownTimeout,
	}
}

func RunServer(ctx context.Context, cfg *Config, handler http.Handler) error {
	log := slog.Default().With("alias", cfg.Alias)

	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("listen http: %w", err)
	}

	httpSrv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	errChan := make(chan error, 1)
	go func() {
		log.InfoContext(ctx, "starting http server", "address", cfg.Addr)
		serveErr := httpSrv.Serve(lis)
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errChan <- fmt.Errorf("serve http server: %w", serveErr)
		}
	}()

	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		err = shutdown(httpSrv, cfg.ShutdownTimeout)
		if err != nil {
			log.ErrorContext(ctx, "shutting down http server", slog.String("error", err.Error()))
		}

		log.InfoContext(ctx, "http server stopped gracefully")
		return nil
	}
}

func shutdown(httpSrv *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return httpSrv.Shutdown(ctx)
}
