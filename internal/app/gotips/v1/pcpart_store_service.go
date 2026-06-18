package v1

import (
	"context"
	"fmt"

	pb "github.com/KlementevTech/gotips/api/gen/go/gotips/v1"
	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/KlementevTech/gotips/internal/domain/repository"
	"github.com/google/uuid"
)

type PCPartStoreService struct {
	pb.UnimplementedPcPartStoreServiceServer

	repo repository.PCPartRepository
}

func NewPCPartStoreService(repo repository.PCPartRepository) *PCPartStoreService {
	return &PCPartStoreService{
		repo: repo,
	}
}

func (s *PCPartStoreService) CreatePcPart(
	ctx context.Context,
	req *pb.CreatePcPartRequest,
) (*pb.CreatePcPartResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s, %w", req.GetId(), domain.ErrInvalidArgument)
	}

	if req.GetName() == "" {
		return nil, fmt.Errorf("invalid name: %s, %w", req.GetName(), domain.ErrInvalidArgument)
	}

	part, err := s.repo.CreatePcPart(ctx, repository.CreatePcPartParams{
		ID:   id,
		Name: req.GetName(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreatePcPartResponse{
		PcPart: toPcPartPb(part),
	}, nil
}

func (s *PCPartStoreService) GetPcPart(
	ctx context.Context,
	req *pb.GetPcPartRequest,
) (*pb.GetPcPartResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s, %w", req.GetId(), domain.ErrInvalidArgument)
	}

	part, err := s.repo.GetPcPartByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.GetPcPartResponse{
		PcPart: toPcPartPb(part),
	}, nil
}

func (s *PCPartStoreService) UpdatePcPart(
	ctx context.Context,
	req *pb.UpdatePcPartRequest,
) (*pb.UpdatePcPartResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s, %w", req.GetId(), domain.ErrInvalidArgument)
	}

	version := fromVersionPb(req.GetVersion())

	params := repository.UpdatePcPartParams{
		ID:      id,
		Version: version,
	}

	for _, p := range req.GetUpdateMask().GetPaths() {
		if p == "name" {
			params.Name = req.GetName()
		}
	}

	part, err := s.repo.UpdatePcPart(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePcPartResponse{
		PcPart: toPcPartPb(part),
	}, nil
}

func (s *PCPartStoreService) DeletePcPart(
	ctx context.Context,
	req *pb.DeletePcPartRequest,
) (*pb.DeletePcPartResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s, %w", req.GetId(), domain.ErrInvalidArgument)
	}

	part, err := s.repo.SoftDeletePcPart(ctx, id, int(req.GetVersion()))
	if err != nil {
		return nil, err
	}

	return &pb.DeletePcPartResponse{
		PcPart: toPcPartPb(part),
	}, nil
}

func (s *PCPartStoreService) GetPcPartsRecent(
	ctx context.Context,
	req *pb.GetPcPartsRecentRequest,
) (*pb.GetPcPartsRecentResponse, error) {
	limit := int32(req.GetLimit().Number())
	if limit == 0 {
		return nil, fmt.Errorf("invalid limit: %d, %w", limit, domain.ErrInvalidArgument)
	}

	res, err := s.repo.GetPcPartsRecent(ctx, limit)
	if err != nil {
		return nil, err
	}

	return &pb.GetPcPartsRecentResponse{
		Items: toPcPartsPb(res),
	}, nil
}
