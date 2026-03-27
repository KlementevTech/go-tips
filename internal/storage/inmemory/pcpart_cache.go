package inmemory

import (
	"context"
	"time"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/hashicorp/golang-lru/v2/simplelru"
)

type pcPartStore interface {
	Create(context.Context, *domain.PcPart) error
	GetByID(context.Context, uuid.UUID) (*domain.PcPart, error)
	Update(context.Context, *domain.PcPart) error
}

type cacheKey struct {
	id uuid.UUID
}

type PcPartCache struct {
	inner simplelru.LRUCache[cacheKey, *domain.PcPart]
	store pcPartStore
}

func NewPcPartCache(store pcPartStore, size int, ttl time.Duration) *PcPartCache {
	inner := expirable.NewLRU[cacheKey, *domain.PcPart](size, nil, ttl)

	return &PcPartCache{
		store: store,
		inner: inner,
	}
}

func (c *PcPartCache) Create(ctx context.Context, mdl *domain.PcPart) error {
	return c.store.Create(ctx, mdl)
}

func (c *PcPartCache) GetByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error) {
	key := cacheKey{id: id}

	if v, ok := c.inner.Get(key); ok {
		return v, nil
	}

	v, err := c.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	c.inner.Add(key, v)
	return v, nil
}

func (c *PcPartCache) Update(ctx context.Context, mdl *domain.PcPart) error {
	err := c.store.Update(ctx, mdl)
	if err != nil {
		return err
	}

	key := cacheKey{id: mdl.ID}
	c.inner.Remove(key)

	return nil
}
