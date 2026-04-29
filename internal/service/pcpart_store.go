package service

import (
	"context"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/google/uuid"
)

type pcPartRepository interface {
	CreatePcPart(ctx context.Context, part *domain.PcPart) error
	GetPcPartByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error)
	GetPcPartsRecent(ctx context.Context, limit int32) ([]*domain.PcPart, error)
	UpdatePcPart(ctx context.Context, part *domain.PcPart) error
	SoftDeletePcPart(ctx context.Context, part *domain.PcPart) error
}

type PcPartStoreService struct {
	repo pcPartRepository
}

func NewPcPartStoreService(repo pcPartRepository) *PcPartStoreService {
	return &PcPartStoreService{
		repo: repo,
	}
}

type CreatePcPartParams struct {
	ID   uuid.UUID
	Name string
}

func (s *PcPartStoreService) Create(ctx context.Context, params *CreatePcPartParams) (*domain.PcPart, error) {
	part := domain.CreatePcPart(params.ID, params.Name)

	err := s.repo.CreatePcPart(ctx, part)
	if err != nil {
		return nil, err
	}
	return part, nil
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
	part, err := s.repo.GetPcPartByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if part.VersionConflict(version) {
		return part, domain.ErrPreconditionFailed
	}

	part.Rename(fields.Name)

	err = s.repo.UpdatePcPart(ctx, part)
	if err != nil {
		return nil, err
	}
	return part, nil
}

func (s *PcPartStoreService) SoftDelete(ctx context.Context, id uuid.UUID, version int) (*domain.PcPart, error) {
	part, err := s.repo.GetPcPartByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if part.VersionConflict(version) {
		return part, domain.ErrPreconditionFailed
	}

	part.MarkAsDeleted()

	err = s.repo.SoftDeletePcPart(ctx, part)
	if err != nil {
		return nil, err
	}
	return part, nil
}
