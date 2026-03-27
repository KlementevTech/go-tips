package domain

import (
	"time"

	"github.com/google/uuid"
)

type PcPart struct {
	ID        uuid.UUID
	Name      string
	Version   int
	CreatedAt time.Time
	DeletedAt *time.Time
}

type NewPcPartFields struct {
	ID   uuid.UUID
	Name string
}

func NewPcPart(fields NewPcPartFields) *PcPart {
	return &PcPart{
		ID:        fields.ID,
		Name:      fields.Name,
		CreatedAt: time.Now().UTC(),
	}
}

func (m *PcPart) Rename(name string) {
	m.Name = name
}

func (m *PcPart) MarkAsDeleted() {
	if m.IsDeleted() {
		return
	}

	now := time.Now().UTC()
	m.DeletedAt = new(now)
}

func (m *PcPart) IsDeleted() bool {
	return m.DeletedAt != nil
}
