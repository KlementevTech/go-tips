package sqlite

import (
	"database/sql"
	"time"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/KlementevTech/gotips/internal/storage/sqlite/sqlc"
)

func toPcPartModel(row sqlc.PcPart) *domain.PcPart {
	return &domain.PcPart{
		ID:        row.ID,
		Name:      row.Name,
		Version:   toVersion(row.Version),
		CreatedAt: row.CreatedAt,
		DeletedAt: toTimePtr(row.DeletedAt),
	}
}

func toTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func fromTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

func fromVersion(v int) int64 {
	return int64(v)
}

func toVersion(v int64) int {
	return int(v)
}
