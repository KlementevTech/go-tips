package internal

import (
	"context"
	"flag"
	"fmt"
	"syscall"

	"github.com/KlementevTech/gotips/internal/config"
	catalogv1 "github.com/KlementevTech/gotips/internal/grpc/handlers/catalog/v1"
	"github.com/KlementevTech/gotips/internal/pprof"
	"github.com/KlementevTech/gotips/internal/service"
	"github.com/KlementevTech/gotips/internal/storage/inmemory"
	"github.com/KlementevTech/gotips/internal/storage/sqlite"
	"github.com/KlementevTech/gotips/internal/transport/grpc"
	"golang.org/x/sync/errgroup"
)

var (
	app     = "go-tips"
	version = "unknown"
)

func Run() error {
	var cfgPath string
	flag.StringVar(&cfgPath, "c", "", "config file path")
	flag.Parse()

	return run(cfgPath)
}

func run(cfgPath string) error {
	changeLvl := defaultJSONLogger(app, version)

	cfg, err := config.LoadFromFile(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	err = changeLvl(cfg.Logger.Level)
	if err != nil {
		return fmt.Errorf("failed to change log level: %w", err)
	}

	ctx, cancel := waitForSignal(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, closeDB, err := sqlite.NewDB(ctx, cfg.SQLite)
	if err != nil {
		return fmt.Errorf("init sqlite: %w", err)
	}
	defer closeDB()

	pcPartsRepo := inmemory.NewPcPartCache(sqlite.NewPcPartStorage(db), cfg.Cache.Size, cfg.Cache.TTL)
	pcPartService := service.NewPcPartService(pcPartsRepo)
	pcPartHandler := catalogv1.NewPcPartStoreHandler(pcPartService)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return grpc.RunServer(
			gCtx,
			cfg.GRPC,
			catalogv1.RegisterHandlersFunc(pcPartHandler),
		)
	})

	if cfg.Pprof.Enable {
		g.Go(func() error {
			return pprof.RunServer(gCtx, cfg.Pprof)
		})
	}

	return g.Wait()
}
