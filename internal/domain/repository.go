package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreatePcPart(ctx context.Context, params CreatePcPartParams) (*PcPart, error)
	GetPcPartByID(ctx context.Context, id uuid.UUID) (*PcPart, error)
	UpdatePcPart(ctx context.Context, params UpdatePcPartParams) (*PcPart, error)
	GetPcPartsRecent(ctx context.Context, limit int32) ([]*PcPart, error)
	SoftDeletePcPart(ctx context.Context, id uuid.UUID, version int) (*PcPart, error)
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
