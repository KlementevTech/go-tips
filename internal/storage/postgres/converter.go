package postgres

import (
	"time"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/KlementevTech/gotips/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func toPcPart(row sqlc.PcPart) *domain.PcPart {
	return &domain.PcPart{
		ID:        row.ID,
		IDString:  row.ID.String(),
		Name:      row.Name,
		Version:   int(row.Version),
		CreatedAt: row.CreatedAt,
		DeletedAt: toTimePtr(row.DeletedAt),
	}
}

func toPcParts(rows []sqlc.PcPart) []*domain.PcPart {
	if len(rows) == 0 {
		return nil
	}

	result := make([]*domain.PcPart, len(rows))
	for i, row := range rows {
		result[i] = toPcPart(row)
	}
	return result
}

func toTimePtr(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
