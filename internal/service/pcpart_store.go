package service

import (
	"context"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/google/uuid"
)

type PcPartStoreService struct {
	repo domain.Repository
}

func NewPcPartStoreService(repo domain.Repository) *PcPartStoreService {
	return &PcPartStoreService{
		repo: repo,
	}
}

type CreatePcPartParams struct {
	ID   uuid.UUID
	Name string
}

func (s *PcPartStoreService) Create(ctx context.Context, params *CreatePcPartParams) (*domain.PcPart, error) {
	return s.repo.CreatePcPart(ctx, domain.CreatePcPartParams{
		ID:   params.ID,
		Name: params.Name,
	})
}

func (s *PcPartStoreService) GetByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error) {
	return s.repo.GetPcPartByID(ctx, id)
}

func (s *PcPartStoreService) GetPcPartsRecent(ctx context.Context, limit int32) ([]*domain.PcPart, error) {
	return s.repo.GetPcPartsRecent(ctx, limit)
}

type UpdatePcPartFields struct {
	Name string
}

func (s *PcPartStoreService) Update(
	ctx context.Context,
	id uuid.UUID,
	version int,
	fields UpdatePcPartFields,
) (*domain.PcPart, error) {
	return s.repo.UpdatePcPart(ctx, domain.UpdatePcPartParams{
		ID:      id,
		Version: version,
		Name:    fields.Name,
	})
}

func (s *PcPartStoreService) SoftDelete(ctx context.Context, id uuid.UUID, version int) (*domain.PcPart, error) {
	return s.repo.SoftDeletePcPart(ctx, id, version)
}
