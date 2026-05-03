package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/KlementevTech/gotips/internal/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
)

var (
	seedsCount    = 1_000
	maxSeedsCount = 100_000
	chanLen       = 100
	minVersion    = 1
	maxVersion    = 5
)

func addSeedsUpCmd(root *cobra.Command) {
	seedsUpCmd := &cobra.Command{
		Use:   "seeds-up",
		Short: "Seeds up",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			n, _ := cmd.Flags().GetInt("count")
			return runSeedsUp(ctx, cfg, n)
		},
	}

	seedsUpCmd.Flags().Int("count", seedsCount, "seeds count")

	root.AddCommand(seedsUpCmd)
}

func runSeedsUp(ctx context.Context, cfg *config.Config, count int) error {
	if err := validateCount(count); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := pgx.Connect(ctx, cfg.Postgres.DSN)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	fields := []string{"id", "name", "version", "created_at"}

	type Seed struct {
		id        uuid.UUID
		name      string
		version   int
		createdAt time.Time
	}

	seedsChan := make(chan Seed, chanLen)

	slog.Default().InfoContext(ctx, "starting seeder", "count", count)
	go func() {
		defer close(seedsChan)
		var seed Seed
		for range count {
			id, _ := uuid.NewV7()

			seed.id = id
			seed.version = gofakeit.Number(minVersion, maxVersion)
			seed.name = fmt.Sprintf("%s %s %d", gofakeit.Company(), gofakeit.Adjective(), seed.version)
			seed.createdAt = time.Now().UTC()

			select {
			case seedsChan <- seed:
			case <-ctx.Done():
				return
			}
		}
	}()

	copyCount, since, err := func() (int64, time.Duration, error) {
		slog.Default().InfoContext(ctx, "starting batch insert")
		start := time.Now()

		copyCount, copyErr := conn.CopyFrom(
			ctx,
			pgx.Identifier{"pc_parts"},
			fields,
			pgx.CopyFromFunc(func() ([]any, error) {
				seed, ok := <-seedsChan
				if !ok {
					return nil, nil
				}

				return []any{
					seed.id,
					seed.name,
					seed.version,
					seed.createdAt,
				}, nil
			}),
		)
		if copyErr != nil {
			return 0, time.Duration(0), fmt.Errorf("failed to insert data: %w", copyErr)
		}

		return copyCount, time.Since(start), nil
	}()
	if err != nil {
		return err
	}

	slog.Default().InfoContext(ctx, "successfully inserted", "copied", copyCount, "duration", since.String())
	return nil
}

func validateCount(n int) error {
	if n > 0 && n <= maxSeedsCount {
		return nil
	}
	return fmt.Errorf("invalid count: %d", n)
}
