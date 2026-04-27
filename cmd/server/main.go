package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/KlementevTech/gotips/internal"
)

var (
	app     = "go-tips"
	version = "unknown"
	config  string
)

func main() {
	flag.StringVar(&config, "c", "", "config file path")
	flag.Parse()

	err := internal.Run(app, version, config)
	if err != nil {
		slog.Default().Error("failed to run", "error", err)
		os.Exit(1)
	}
}
