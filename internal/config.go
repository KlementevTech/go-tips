package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type GRPCConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	MaxRecvMsgSize    int           `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSize    int           `mapstructure:"max_send_msg_size"`
	MaxConnectionIdle time.Duration `mapstructure:"max_connection_idle"`
	MaxConnectionAge  time.Duration `mapstructure:"max_connection_age"`
	EnableReflection  bool          `mapstructure:"enable_reflection"`
	EnableHealth      bool          `mapstructure:"enable_health"`
}

type PprofConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	Enable            bool          `mapstructure:"enable"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
}

type LogConfig struct {
	Format string `mapstructure:"format"`
	Level  string `mapstructure:"level"`
}

type Config struct {
	GRPC  GRPCConfig  `mapstructure:"grpc"`
	Pprof PprofConfig `mapstructure:"pprof"`
	Log   LogConfig   `mapstructure:"log"`
}

func LoadConfig(path string) (*Config, error) {
	return loadConfig(path, "")
}

func loadConfig(path, envPrefix string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config file path is required")
	}

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
