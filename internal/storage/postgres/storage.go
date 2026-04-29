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

func (s *Storage) CreatePcPart(ctx context.Context, part *domain.PcPart) error {
	const op = "postgres.CreatePcPart"

	row, err := s.db(ctx).CreatePcPart(ctx, sqlc.CreatePcPartParams{
		ID:        part.ID,
		Name:      part.Name,
		Version:   1,
		CreatedAt: part.CreatedAt,
	})
	if err != nil {
		if isAlreadyExists(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	*part = *toPcPart(row)
	return nil
}

const pgAlreadyExistsCode = "23505"

func isAlreadyExists(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == pgAlreadyExistsCode
	}
	return false
}

func (s *Storage) GetPcPartByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error) {
	const op = "postgres.GetPcPartsRecent"

	row, err := s.db(ctx).GetPcPart(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return toPcPart(row), nil
}

func (s *Storage) UpdatePcPart(ctx context.Context, part *domain.PcPart) error {
	const op = "postgres.UpdatePcPart"

	row, err := s.db(ctx).UpdatePcPart(ctx, sqlc.UpdatePcPartParams{
		ID:         part.ID,
		Name:       part.Name,
		Version:    int64(part.Version + 1),
		OldVersion: int64(part.Version),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrPreconditionFailed
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	*part = *toPcPart(row)
	return nil
}

func (s *Storage) GetPcPartsRecent(ctx context.Context, limit int32) ([]*domain.PcPart, error) {
	const op = "postgres.GetPcPartsRecent"

	rows, err := s.db(ctx).GetPcPartsRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return toPcParts(rows), nil
}

func (s *Storage) SoftDeletePcPart(ctx context.Context, part *domain.PcPart) error {
	const op = "postgres.SoftDeletePcPart"

	row, err := s.db(ctx).SoftDeletePcPart(ctx, sqlc.SoftDeletePcPartParams{
		ID:         part.ID,
		DeletedAt:  fromTimePtr(part.DeletedAt),
		Version:    int64(part.Version + 1),
		OldVersion: int64(part.Version),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrPreconditionFailed
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	*part = *toPcPart(row)
	return nil
}

func (s *Storage) db(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(txCtxKey{}).(pgx.Tx); ok {
		return s.queries.WithTx(tx)
	}
	return s.queries
}
