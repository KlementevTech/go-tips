package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/KlementevTech/gotips/internal/storage/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	queries *sqlc.Queries
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		queries: sqlc.New(pool),
	}
}

func (s *Storage) CreatePcPart(ctx context.Context, params domain.CreatePcPartParams) (*domain.PcPart, error) {
	const op = "postgres.CreatePcPart"

	row, err := s.db(ctx).CreatePcPart(ctx, sqlc.CreatePcPartParams{
		ID:   params.ID,
		Name: params.Name,
	})
	if err != nil {
		if isAlreadyExists(err) {
			return nil, domain.ErrAlreadyExists
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return toPcPart(row), nil
}

const pgAlreadyExistsCode = "23505"

func isAlreadyExists(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == pgAlreadyExistsCode
	}
	return false
}

func (s *Storage) GetPcPartByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error) {
	const op = "postgres.GetPcPartByID"

	row, err := s.db(ctx).GetPcPart(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return toPcPart(row), nil
}

func (s *Storage) UpdatePcPart(ctx context.Context, params domain.UpdatePcPartParams) (*domain.PcPart, error) {
	const op = "postgres.UpdatePcPart"

	row, err := s.db(ctx).UpdatePcPart(ctx, sqlc.UpdatePcPartParams{
		ID:      params.ID,
		Version: int64(params.Version),
		Name:    params.Name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPreconditionFailed
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return toPcPart(row), nil
}

func (s *Storage) GetPcPartsRecent(ctx context.Context, limit int32) ([]*domain.PcPart, error) {
	const op = "postgres.GetPcPartsRecent"

	rows, err := s.db(ctx).GetPcPartsRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return toPcParts(rows), nil
}

func (s *Storage) SoftDeletePcPart(ctx context.Context, id uuid.UUID, version int) (*domain.PcPart, error) {
	const op = "postgres.SoftDeletePcPart"

	row, err := s.db(ctx).SoftDeletePcPart(ctx, sqlc.SoftDeletePcPartParams{
		ID:      id,
		Version: int64(version),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPreconditionFailed
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return toPcPart(row), nil
}

func (s *Storage) db(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(txCtxKey{}).(pgx.Tx); ok {
		return s.queries.WithTx(tx)
	}
	return s.queries
}
