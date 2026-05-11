package main

import (
	"log/slog"
	"os"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		slog.Default().Error("failed to execute root command", "error", err)
		os.Exit(1)
	}
}
