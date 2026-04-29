package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/KlementevTech/gotips/internal"
)

var version = "dev"

func main() {
	var config string
	flag.StringVar(&config, "c", "", "config file path")
	flag.Parse()

	err := internal.Run(version, config)
	if err != nil {
		slog.Default().Error("failed to run", "error", err)
		os.Exit(1)
	}
}
