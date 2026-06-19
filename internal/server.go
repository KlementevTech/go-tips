package internal

import (
	"context"
	"fmt"
	"syscall"

	gotipsv1 "github.com/KlementevTech/gotips/internal/app/gotips/v1"
	"github.com/KlementevTech/gotips/internal/app/todos"
	"github.com/KlementevTech/gotips/internal/grpc"
	"github.com/KlementevTech/gotips/internal/http"
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

	ctx, cancel := waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return grpc.RunServer(
			ctx,
			gotipsv1.NewPCPartStoreService(pcPartCache),
			todos.NewService(),
			cfg.GRPC,
		)
	})

	if cfg.Pprof.Enable {
		g.Go(func() error {
			return http.RunServer(
				ctx,
				cfg.Pprof.Address,
				pprof.RegisterRoutes(),
				http.WithAlias("pprof"),
			)
		})
	}

	return g.Wait()
}
