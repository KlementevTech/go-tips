package log

import (
	"fmt"
	"log/slog"
	"os"
)

type SetLevelFunc func(lvl string) error

func SetupJSONLog(opts ...Option) SetLevelFunc {
	lvl := new(slog.LevelVar)

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: lvl,
		},
	)

	log := slog.New(handler)

	for _, opt := range opts {
		log = opt(log)
	}

	slog.SetDefault(log)

	return func(raw string) error {
		if err := lvl.UnmarshalText([]byte(raw)); err != nil {
			return fmt.Errorf("changing log level: %w", err)
		}
		return nil
	}
}

type Option func(*slog.Logger) *slog.Logger

func WithVersion(version string) Option {
	return func(l *slog.Logger) *slog.Logger {
		return l.With("version", version)
	}
}
