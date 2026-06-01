package repository

import (
	"context"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/google/uuid"
)

type PCPartRepository interface {
	CreatePcPart(ctx context.Context, params CreatePcPartParams) (*domain.PcPart, error)
	GetPcPartByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error)
	UpdatePcPart(ctx context.Context, params UpdatePcPartParams) (*domain.PcPart, error)
	GetPcPartsRecent(ctx context.Context, limit int32) ([]*domain.PcPart, error)
	SoftDeletePcPart(ctx context.Context, id uuid.UUID, version int) (*domain.PcPart, error)
}

type CreatePcPartParams struct {
	ID   uuid.UUID
	Name string
}

type UpdatePcPartParams struct {
	ID      uuid.UUID
	Version int
	Name    string
}
