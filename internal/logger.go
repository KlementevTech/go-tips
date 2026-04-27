package internal

import (
	"fmt"
	"log/slog"
	"os"
)

type changeLogLevel func(lvl string) error

func defaultJSONLogger(app, version string) changeLogLevel {
	vl := new(slog.LevelVar)
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: vl,
		},
	)
	log := slog.New(handler).With(
		slog.String("app", app),
		slog.String("version", version),
	)

	slog.SetDefault(log)

	return func(lvl string) error {
		err := vl.UnmarshalText([]byte(lvl))
		if err != nil {
			return fmt.Errorf("invalid log level: %w", err)
		}

		slog.Default().Info("log level changed", "level", vl.Level().String())
		return nil
	}
}
