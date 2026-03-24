package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

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

	cfg, err := internal.LoadConfig(path)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	err = internal.InitLog(cfg.Log)
	if err != nil {
		return fmt.Errorf("init log: %w", err)
	}

	ctx := context.Background()
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return internal.RunGRPCServer(gCtx, cfg.GRPC)
	})

	if cfg.Pprof.Enable {
		g.Go(func() error {
			return internal.RunPprofServer(gCtx, cfg.Pprof)
		})
	}

	g.Go(func() error {
		sig := internal.WaitForSignals(gCtx)
		if sig != nil {
			slog.Default().Info("received signal", "signal", sig.String())
			return context.Canceled
		}
		return nil
	})

	err = g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
