package v1

import (
	"time"

	pb "github.com/KlementevTech/gotips/api/gen/pb/catalog/v1"
	"github.com/KlementevTech/gotips/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toPcPartPb(m *domain.PcPart) *pb.PcPart {
	if m == nil {
		return nil
	}

	return &pb.PcPart{
		Id:        m.ID.String(),
		Name:      m.Name,
		Version:   toVersionPb(m.Version),
		CreatedAt: timestamppb.New(m.CreatedAt),
		DeletedAt: toTimestampPb(m.DeletedAt),
	}
}

func toTimestampPb(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func fromVersionPb(version int64) int {
	return int(version)
}

func toVersionPb(version int) int64 {
	return int64(version)
}
