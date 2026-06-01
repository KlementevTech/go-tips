package internal

import (
	"context"
	"fmt"
	"syscall"

	v1 "github.com/KlementevTech/gotips/internal/app/gotips/v1"
	"github.com/KlementevTech/gotips/internal/grpc"
	http2 "github.com/KlementevTech/gotips/internal/http"
	"github.com/KlementevTech/gotips/internal/pprof"
	"github.com/KlementevTech/gotips/internal/storage/cache/pcpart"
	"github.com/KlementevTech/gotips/internal/storage/postgres"
	"golang.org/x/sync/errgroup"
)

func Run(cfg *Config) error {
	ctx := context.Background()

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

	ctx, cancel := waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return grpc.RunServer(ctx, pcPartStoreService, cfg.GRPC)
	})

	if cfg.Pprof.Enable {
		g.Go(func() error {
			return http2.RunServer(
				ctx,
				cfg.Pprof.Address,
				pprof.RegisterRoutes(),
				http2.WithAlias("pprof"),
			)
		})
	}

	return g.Wait()
}
