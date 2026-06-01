package http

import (
	"log/slog"
	"time"
)

type serverOptions struct {
	readHeaderTimeout time.Duration
	shutdownTimeout   time.Duration
	logger            *slog.Logger
}

const (
	shutdownTimeout   = 30 * time.Second
	readHeaderTimeout = 10 * time.Second
)

var defaultServerOptions = serverOptions{
	readHeaderTimeout: readHeaderTimeout,
	shutdownTimeout:   shutdownTimeout,
	logger:            slog.Default(),
}

type ServerOption interface {
	apply(options *serverOptions)
}

type serverOptionFunc func(options *serverOptions)

func (f serverOptionFunc) apply(options *serverOptions) {
	f(options)
}

func WithReadHeaderTimeout(timeout time.Duration) ServerOption {
	return serverOptionFunc(func(options *serverOptions) {
		options.readHeaderTimeout = timeout
	})
}

func WithShutdownTimeout(timeout time.Duration) ServerOption {
	return serverOptionFunc(func(options *serverOptions) {
		options.shutdownTimeout = timeout
	})
}

func WithAlias(alias string) ServerOption {
	return serverOptionFunc(func(options *serverOptions) {
		options.logger = options.logger.With("alias", alias)
	})
}
