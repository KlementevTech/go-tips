package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/google/uuid"
)

type pcPartRepository interface {
	Create(context.Context, *domain.PcPart) error
	GetByID(context.Context, uuid.UUID) (*domain.PcPart, error)
	Update(context.Context, *domain.PcPart) error
}

type PcPartService struct {
	repo pcPartRepository
}

func NewPcPartService(repo pcPartRepository) *PcPartService {
	return &PcPartService{
		repo: repo,
	}
}

type CreatePcPartParams struct {
	ID   uuid.UUID
	Name string
}

func (s *PcPartService) Create(ctx context.Context, params *CreatePcPartParams) (*domain.PcPart, error) {
	fields := domain.NewPcPartFields{
		ID:   params.ID,
		Name: params.Name,
	}
	mdl := domain.NewPcPart(fields)

	err := s.repo.Create(ctx, mdl)
	if err != nil {
		return nil, err
	}
	return mdl, nil
}

func (s *PcPartService) GetByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error) {
	mdl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
		return nil, err
	}
	return mdl, nil
}

type UpdatePcPartParams struct {
	Name string
}

func (s *PcPartService) Update(
	ctx context.Context,
	id uuid.UUID,
	version int,
	params UpdatePcPartParams,
) (*domain.PcPart, error) {
	mdl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if mdl.Version != version {
		return nil, domain.ErrPreconditionFailed
	}

	mdl.Rename(params.Name)

	err = s.repo.Update(ctx, mdl)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to Update pcPart: %w", err)
	}
	return mdl, nil
}

func (s *PcPartService) SoftDelete(ctx context.Context, id uuid.UUID, version int) error {
	mdl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if mdl.Version != version {
		return domain.ErrPreconditionFailed
	}

	mdl.MarkAsDeleted()

	err = s.repo.Update(ctx, mdl)
	if err != nil {
		return err
	}
	return nil
}
