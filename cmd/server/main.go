package main

import (
	"log/slog"
	"os"

	"github.com/KlementevTech/gotips/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		slog.Default().Error("failed to run", "error", err)
		os.Exit(1)
	}
}
