package internal

import (
	"fmt"
	"log/slog"
	"os"
)

func InitLog(cfg LogConfig) error {
	var l slog.Level
	if err := l.UnmarshalText([]byte(cfg.Level)); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	options := &slog.HandlerOptions{
		Level: l,
	}

	var handler slog.Handler
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, options)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, options)
	default:
		return fmt.Errorf("invalid log format: %s", cfg.Format)
	}

	slog.SetDefault(slog.New(handler))
	return nil
}
