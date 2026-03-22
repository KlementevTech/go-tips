package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/KlementevTech/gotips/internal"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		slog.Default().Error("failed to run", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var path string
	flag.StringVar(&path, "c", "", "config file")
	flag.Parse()

	cfg, err := internal.LoadConfigFromFile(path, "")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	notifyCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer func() {
		err = context.Cause(notifyCtx)
		if err != nil && !errors.Is(err, context.Canceled) {
			slog.Default().InfoContext(notifyCtx, "received signal", "signal", err)
		}
	}()

	g, gCtx := errgroup.WithContext(notifyCtx)

	if cfg.Pprof.Enable {
		err = internal.RunPprofServer(gCtx, g, cfg.Pprof)
		if err != nil {
			return fmt.Errorf("run pprof server: %w", err)
		}
	}

	if err = g.Wait(); err != nil {
		return err
	}
	return nil
}
