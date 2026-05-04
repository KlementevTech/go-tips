package v1

import (
	"time"

	pb "github.com/KlementevTech/gotips/api/gen/pb/gotips/v1"
	"github.com/KlementevTech/gotips/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toPcPartPb(m *domain.PcPart) *pb.PcPart {
	if m == nil {
		return nil
	}

	deletedAt := toTimestampPb(m.DeletedAt)

	return &pb.PcPart{
		Id:        m.IDString,
		Name:      m.Name,
		Version:   toVersionPb(m.Version),
		CreatedAt: timestamppb.New(m.CreatedAt),
		DeletedAt: deletedAt,
	}
}

func toPcPartsPb(list []*domain.PcPart) []*pb.PcPart {
	if len(list) == 0 {
		return nil
	}

	result := make([]*pb.PcPart, len(list))
	for i, p := range list {
		result[i] = toPcPartPb(p)
	}
	return result
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
