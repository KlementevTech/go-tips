package internal

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
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

type params struct {
	app     string
	version string
	cfgPath string
}

func Run() error {
	p := params{
		app:     app,
		version: version,
	}

	flag.StringVar(&p.cfgPath, "c", "", "config file path")
	flag.Parse()
	return run(p)
}

func run(p params) error {
	cfg, err := config.LoadFromFile(p.cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	err = setupJSONLogger(p.app, p.version, cfg.Logger.Level)
	if err != nil {
		return fmt.Errorf("setup logger: %w", err)
	}

	slog.Default().Info("running go-tips app", slog.String("config", p.cfgPath))

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
