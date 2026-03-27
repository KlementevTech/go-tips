package v1

import (
	"context"

	pb "github.com/KlementevTech/gotips/api/gen/pb/catalog/v1"
	"github.com/KlementevTech/gotips/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.PcPartServiceServer = (*pcPartStoreHandler)(nil)

type pcPartStoreHandler struct {
	pb.UnimplementedPcPartServiceServer

	service *service.PcPartService
}

func NewPcPartStoreHandler(uc *service.PcPartService) pb.PcPartServiceServer {
	return &pcPartStoreHandler{
		service: uc,
	}
}

func (s *pcPartStoreHandler) CreatePcPart(
	ctx context.Context,
	req *pb.CreatePcPartRequest,
) (*pb.CreatePcPartResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.GetId())
	}

	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %s", req.GetName())
	}

	params := &service.CreatePcPartParams{
		ID:   id,
		Name: req.GetName(),
	}

	mdl, err := s.service.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePcPartResponse{
		PcPart: toPcPartPb(mdl),
	}, nil
}

func (s *pcPartStoreHandler) GetPcPart(
	ctx context.Context,
	req *pb.GetPcPartRequest,
) (*pb.GetPcPartResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.GetId())
	}

	mdl, err := s.service.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.GetPcPartResponse{
		PcPart: toPcPartPb(mdl),
	}, nil
}

func (s *pcPartStoreHandler) UpdatePcPart(
	ctx context.Context,
	req *pb.UpdatePcPartRequest,
) (*pb.UpdatePcPartResponse, error) {
	id, err := uuid.Parse(req.GetPcPart().GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.GetPcPart().GetId())
	}

	if req.GetPcPart().GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %s", req.GetPcPart().GetName())
	}

	version := fromVersionPb(req.GetPcPart().GetVersion())

	var params service.UpdatePcPartParams

	for _, p := range req.GetUpdateMask().GetPaths() {
		if p == "name" {
			params.Name = req.GetPcPart().GetName()
		}
	}

	mdl, err := s.service.Update(ctx, id, version, params)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePcPartResponse{
		PcPart: toPcPartPb(mdl),
	}, nil
}

func (s *pcPartStoreHandler) DeletePcPart(
	ctx context.Context,
	req *pb.DeletePcPartRequest,
) (*pb.DeletePcPartResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.GetId())
	}

	err = s.service.SoftDelete(ctx, id, fromVersionPb(req.GetVersion()))
	if err != nil {
		return nil, err
	}

	return &pb.DeletePcPartResponse{}, nil
}

func RegisterHandlersFunc(handler pb.PcPartServiceServer) func(s *grpc.Server) error {
	return func(s *grpc.Server) error {
		pb.RegisterPcPartServiceServer(s, handler)
		return nil
	}
}
