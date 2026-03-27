package internal

import (
	"fmt"
	"log/slog"
	"os"
)

func setupJSONLogger(app, version, lvl string) error {
	var l slog.Level
	err := l.UnmarshalText([]byte(lvl))
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	options := &slog.HandlerOptions{
		Level: l,
	}

	handler := slog.NewJSONHandler(os.Stdout, options)
	slog.SetDefault(
		slog.New(handler).With(
			slog.String("app", app),
			slog.String("version", version),
		),
	)
	return nil
}
