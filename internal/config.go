package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type PprofConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	Enable            bool          `mapstructure:"enable"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
}

type Config struct {
	Pprof PprofConfig `mapstructure:"pprof"`
}

func LoadConfigFromFile(path, envPrefix string) (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(false)

	viper.SetConfigFile(path)

	var cfg Config
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
