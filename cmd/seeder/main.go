package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/KlementevTech/gotips/pkg/log"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	seedsCount    = 1_000
	maxSeedsCount = 100_000
	seedsChanLen  = 100
	minVersion    = 1
	maxVersion    = 5
)

func main() {
	ctx := context.Background()
	log.SetupJSONLog()
	dsn, ok := os.LookupEnv("POSTGRES_DSN")
	if !ok {
		slog.Default().ErrorContext(ctx, "POSTGRES_DSN env variable not set")
		os.Exit(1)
	}

	var count int
	flag.IntVar(&count, "n", seedsCount, "number of seeds to generate")
	flag.Parse()

	if !isValidCount(count) {
		slog.Default().ErrorContext(ctx, fmt.Sprintf("n must be between 0 and %d", maxSeedsCount))
		os.Exit(1)
	}

	err := run(ctx, dsn, count)
	if err != nil {
		slog.Default().ErrorContext(ctx, "error running seeder", "error", err)
		os.Exit(1)
	}
}

func isValidCount(n int) bool {
	return n > 0 && n <= maxSeedsCount
}

func run(ctx context.Context, dsn string, count int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
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

	seedsChan := make(chan Seed, seedsChanLen)

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
