package internal

import (
	"errors"
	"time"

	"github.com/KlementevTech/gotips/internal/grpc"
	"github.com/KlementevTech/gotips/internal/storage/postgres"
	"github.com/KlementevTech/gotips/pkg/config"
)

type Config struct {
	GRPC     grpc.Config     `mapstructure:"grpc"`
	Postgres postgres.Config `mapstructure:"postgres"`
	Cache    Cache           `mapstructure:"cache"`
	Logger   Logger          `mapstructure:"logger"`
	Pprof    Pprof           `mapstructure:"pprof"`
}

type Logger struct {
	Level string `mapstructure:"level"`
}

type Cache struct {
	Size    int           `mapstructure:"size"`
	TTL     time.Duration `mapstructure:"ttl"`
	Timeout time.Duration `mapstructure:"timeout"`
	Shards  uint32        `mapstructure:"shards"`
}

type Pprof struct {
	Address string `mapstructure:"address"`
	Enable  bool   `mapstructure:"enable"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config file path is required")
	}

	return config.LoadFromFile[Config](path, "")
}
