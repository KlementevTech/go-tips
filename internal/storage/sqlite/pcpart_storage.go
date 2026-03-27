package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/KlementevTech/gotips/internal/storage/sqlite/sqlc"
	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type PcPartStorage struct {
	q   *sqlc.Queries
	sfg singleflight.Group
}

func NewPcPartStorage(db *sql.DB) *PcPartStorage {
	return &PcPartStorage{
		q: sqlc.New(db),
	}
}

func (s *PcPartStorage) Create(ctx context.Context, mdl *domain.PcPart) error {
	row, err := s.queries(ctx).CreatePcPart(ctx, sqlc.CreatePcPartParams{
		ID:        mdl.ID,
		Name:      mdl.Name,
		Version:   1,
		CreatedAt: mdl.CreatedAt,
	})
	if err != nil {
		if isAlreadyExists(err) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("sqlite: failed to call CreatePcPart: %w", err)
	}

	*mdl = *toPcPartModel(row)
	return nil
}

func (s *PcPartStorage) GetByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error) {
	row, err, _ := s.sfg.Do(id.String(), func() (any, error) {
		row, err := s.queries(ctx).GetPcPart(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrNotFound
			}
			return nil, fmt.Errorf("sqlite: failed to call GetPcPart: %w", err)
		}
		return row, nil
	})

	if err != nil {
		return nil, err
	}

	//nolint: errcheck // nolint
	return toPcPartModel(row.(sqlc.PcPart)), nil
}

func (s *PcPartStorage) Update(ctx context.Context, mdl *domain.PcPart) error {
	row, err := s.queries(ctx).UpdatePcPart(ctx, sqlc.UpdatePcPartParams{
		ID:         mdl.ID,
		Name:       mdl.Name,
		Version:    fromVersion(mdl.Version + 1),
		OldVersion: fromVersion(mdl.Version),
		DeletedAt:  fromTimePtr(mdl.DeletedAt),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrPreconditionFailed
		}
		return fmt.Errorf("sqlite: failed to call UpdatePcPart: %w", err)
	}

	*mdl = *toPcPartModel(row)
	return nil
}

func (s *PcPartStorage) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(txCtxKey{}).(*sql.Tx); ok {
		return s.q.WithTx(tx)
	}
	return s.q
}

const (
	constraintPrimaryKeyCode = sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY
	constraintUniqueCode     = sqlite3.SQLITE_CONSTRAINT_UNIQUE
)

func isAlreadyExists(err error) bool {
	if asErr, ok := errors.AsType[*sqlite.Error](err); ok {
		code := asErr.Code()
		if code == constraintPrimaryKeyCode || code == constraintUniqueCode {
			return true
		}
	}
	return false
}
