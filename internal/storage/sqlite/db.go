package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"
)

type Config struct {
	DBFile          string        `mapstructure:"db_file"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnectTimeout  time.Duration `mapstructure:"connect_timeout"`
}

type CloseFunc func()

func NewDB(ctx context.Context, cfg Config) (*sql.DB, CloseFunc, error) {
	absPath, err := filepath.Abs(cfg.DBFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL", absPath)

	slog.Default().InfoContext(ctx, "opening sqlite", slog.String("dsn", dsn))
	var db *sql.DB
	db, err = sql.Open("sqlite", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid dsn: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("database unreachable (timeout %v): %w", cfg.ConnectTimeout, err)
	}

	return db, func() {
		if err = db.Close(); err != nil {
			slog.Default().ErrorContext(ctx, "failed to close sqlite", "error", err)
		} else {
			slog.Default().InfoContext(ctx, "sqlite closed gracefully")
		}
	}, nil
}
