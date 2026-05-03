package domain

import (
	"time"

	"github.com/google/uuid"
)

type PcPart struct {
	ID        uuid.UUID
	IDString  string
	Name      string
	Version   int
	CreatedAt time.Time
	DeletedAt *time.Time
}
