package v1

import (
	"context"

	pb "github.com/KlementevTech/gotips/api/gen/pb/gotips/v1"
	"github.com/KlementevTech/gotips/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.PcPartStoreServiceServer = (*pcPartStoreHandler)(nil)

type pcPartStoreHandler struct {
	pb.UnimplementedPcPartStoreServiceServer

	pcPartsStoreService *service.PcPartStoreService
}

func NewPcPartStoreHandler(
	pcPartsStoreService *service.PcPartStoreService,
) pb.PcPartStoreServiceServer {
	return &pcPartStoreHandler{
		pcPartsStoreService: pcPartsStoreService,
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

	part, err := s.pcPartsStoreService.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePcPartResponse{
		PcPart: toPcPartPb(part),
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

	part, err := s.pcPartsStoreService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.GetPcPartResponse{
		PcPart: toPcPartPb(part),
	}, nil
}

func (s *pcPartStoreHandler) UpdatePcPart(
	ctx context.Context,
	req *pb.UpdatePcPartRequest,
) (*pb.UpdatePcPartResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", req.GetId())
	}

	version := fromVersionPb(req.GetVersion())

	var params service.UpdatePcPartFields

	for _, p := range req.GetUpdateMask().GetPaths() {
		if p == "name" {
			params.Name = req.GetName()
		}
	}

	part, err := s.pcPartsStoreService.Update(ctx, id, version, params)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePcPartResponse{
		PcPart: toPcPartPb(part),
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

	part, err := s.pcPartsStoreService.SoftDelete(ctx, id, fromVersionPb(req.GetVersion()))
	if err != nil {
		return nil, err
	}

	return &pb.DeletePcPartResponse{
		PcPart: toPcPartPb(part),
	}, nil
}

func (s *pcPartStoreHandler) GetPcPartsRecent(
	ctx context.Context,
	req *pb.GetPcPartsRecentRequest,
) (*pb.GetPcPartsRecentResponse, error) {
	limit := int32(req.GetLimit().Number())
	if limit == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid limit: %d", req.GetLimit().Number())
	}

	res, err := s.pcPartsStoreService.GetPcPartsRecent(ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get pc parts recent: %s", err)
	}

	return &pb.GetPcPartsRecentResponse{
		Items: toPcPartsPb(res),
	}, nil
}

func RegisterHandlersFunc(handler pb.PcPartStoreServiceServer) func(s *grpc.Server) error {
	return func(s *grpc.Server) error {
		pb.RegisterPcPartStoreServiceServer(s, handler)
		return nil
	}
}
