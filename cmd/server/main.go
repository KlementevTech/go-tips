package main

import (
	"context"
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

	cfg, err := internal.LoadConfig(path)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	err = internal.InitLog(cfg.Log)
	if err != nil {
		return fmt.Errorf("init log: %w", err)
	}

	g, gCtx := errgroup.WithContext(context.Background())
	ctx, stop := withInterrupt(gCtx)
	defer stop()

	if cfg.Pprof.Enable {
		err = internal.RunPprofServer(ctx, g, cfg.Pprof)
		if err != nil {
			return fmt.Errorf("run pprof server: %w", err)
		}
	}

	return g.Wait()
}

func withInterrupt(ctx context.Context) (context.Context, context.CancelFunc) {
	newCtx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer signal.Stop(sigs)

		select {
		case sig := <-sigs:
			slog.Default().Info("received signal", "signal", sig.String())
			cancel()
		case <-newCtx.Done():
		}
	}()

	return newCtx, cancel
}
