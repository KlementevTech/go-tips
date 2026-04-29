package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/KlementevTech/gotips/internal/pprof"
	"github.com/KlementevTech/gotips/internal/storage/sqlite"
	"github.com/KlementevTech/gotips/internal/transport/grpc"
	"github.com/spf13/viper"
)

type Logger struct {
	Level string `mapstructure:"level"`
}

type Cache struct {
	Size int           `mapstructure:"size"`
	TTL  time.Duration `mapstructure:"ttl"`
}

type Config struct {
	GRPC   grpc.Config   `mapstructure:"grpc"`
	Pprof  pprof.Config  `mapstructure:"pprof"`
	SQLite sqlite.Config `mapstructure:"sqlite"`
	Cache  Cache         `mapstructure:"cache"`
	Logger Logger        `mapstructure:"logger"`
}

func LoadFromFile(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config file path is required")
	}

	cfg, err := loadFromFile[Config](path, "")
	if err != nil {
		return nil, err
	}

	slog.Default().Info("config loaded", "config", path)
	return cfg, nil
}

func loadFromFile[T any](path, envPrefix string) (*T, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(false)

	viper.SetConfigFile(path)

	var cfg T
	if err := viper.ReadInConfig(); err != nil {
		if tErr, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			return nil, fmt.Errorf("config file not found: %w", tErr)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}
	return &cfg, nil
}
