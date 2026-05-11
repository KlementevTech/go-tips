package internal

import (
	"context"
	"fmt"
	"log/slog"
	"syscall"

	v1 "github.com/KlementevTech/gotips/internal/app/gotips/v1"
	"github.com/KlementevTech/gotips/internal/config"
	"github.com/KlementevTech/gotips/internal/pprof"
	"github.com/KlementevTech/gotips/internal/storage/cache/pcpart"
	"github.com/KlementevTech/gotips/internal/storage/postgres"
	"github.com/KlementevTech/gotips/internal/transport/grpc"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pgPool, closePgPool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("create postgres pool: %w", err)
	}
	defer closePgPool()

	pgStorage := postgres.NewStorage(pgPool)
	pcPartCache := pcpart.NewLRUCache(pgStorage, &pcpart.LRUCacheConfig{
		Size:    cfg.Cache.Size,
		TTL:     cfg.Cache.TTL,
		Timeout: cfg.Cache.Timeout,
		Shards:  cfg.Cache.Shards,
	})

	pcPartStoreService := v1.NewPCPartStoreService(pcPartCache)

	slog.Default().InfoContext(ctx, "service initialized, starting servers")

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return grpc.RunServer(gCtx, pcPartStoreService, cfg.GRPC)
	})

	if cfg.Pprof.Enable {
		g.Go(func() error {
			return pprof.RunServer(gCtx, cfg.Pprof)
		})
	}

	return g.Wait()
}
