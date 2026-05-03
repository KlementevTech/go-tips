package main

import (
	"fmt"
	"log/slog"

	"github.com/KlementevTech/gotips/internal/config"
	"github.com/KlementevTech/gotips/pkg/log"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	cfgFile string
	cfg     *config.Config
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gotips",
		Short: "Go-Tips Service",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			setLevel := log.SetupJSONLog(log.WithVersion(version))

			var err error
			cfg, err = config.LoadFromFile(cfgFile)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			slog.Default().Info("loaded config", "file", cfgFile)

			err = setLevel(cfg.Logger.Level)
			if err != nil {
				return fmt.Errorf("failed to set log level: %w", err)
			}

			return nil
		},
	}

	root.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"config/config.local.toml",
		"path to config file",
	)

	addServeCmd(root)
	addSeedsUpCmd(root)

	return root
}
