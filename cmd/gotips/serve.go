package main

import (
	"github.com/KlementevTech/gotips/internal"
	"github.com/spf13/cobra"
)

func addServeCmd(root *cobra.Command) {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve gRPC, Pprof servers",
		RunE: func(_ *cobra.Command, _ []string) error {
			return internal.Run(cfg)
		},
	}

	root.AddCommand(serveCmd)
}
