package main

import (
	"context"

	"github.com/KlementevTech/gotips/internal"
	"github.com/spf13/cobra"
)

func addServeCmd(root *cobra.Command) {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve gRPC, Pprof servers",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			return internal.Run(ctx, cfg)
		},
	}

	root.AddCommand(serveCmd)
}
