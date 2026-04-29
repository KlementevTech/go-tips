package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN               string        `mapstructure:"dsn"`
	ConnTimeout       time.Duration `mapstructure:"conn_timeout"`
	MaxConns          int32         `mapstructure:"max_conns"`
	MinIdleConns      int32         `mapstructure:"min_idle_conns"`
	MaxConnLifetime   time.Duration `mapstructure:"max_conn_lifetime"`
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period"`
	AppName           string        `mapstructure:"app_name"`
}

func NewPool(ctx context.Context, cfg *Config) (*pgxpool.Pool, func(), error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, nil, err
	}

	dbConfig.ConnConfig.ConnectTimeout = cfg.ConnTimeout                // 5 * time.Second
	dbConfig.ConnConfig.RuntimeParams["application_name"] = cfg.AppName // gotips
	dbConfig.MaxConns = cfg.MaxConns                                    // 25
	dbConfig.MinIdleConns = cfg.MinIdleConns                            // 5
	dbConfig.MaxConnLifetime = cfg.MaxConnLifetime                      // 30 * time.Minute
	dbConfig.HealthCheckPeriod = cfg.HealthCheckPeriod                  // 1 * time.Minute

	/*
		// Если вы используете кастомные типы или логирование сессий
		dbConfig.AfterConnect = func( context.Context, conn *pgx.Conn) error {
			// Настройка конкретной "трубки"
			connString := conn.Config().ConnString()
			slog.Default().InfoContext(ctx, "postgres connected", "dsn", connString)
			return nil
		}
	*/

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to ping postgres database: %w", err)
	}

	return pool, func() {
		pool.Close()
		slog.Default().InfoContext(ctx, "Postgres closed gracefully")
	}, nil
}
