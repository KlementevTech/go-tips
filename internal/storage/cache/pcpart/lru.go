package pcpart

import (
	"context"
	"fmt"
	"time"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"golang.org/x/sync/singleflight"
)

type pcPartRepository interface {
	CreatePcPart(ctx context.Context, part *domain.PcPart) error
	GetPcPartByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error)
	GetPcPartsRecent(ctx context.Context, limit int32) ([]*domain.PcPart, error)
	UpdatePcPart(ctx context.Context, part *domain.PcPart) error
	SoftDeletePcPart(ctx context.Context, part *domain.PcPart) error
}

type LRUCache struct {
	storage pcPartRepository
	recent  *expirable.LRU[string, []*domain.PcPart]
	inner   *expirable.LRU[string, *domain.PcPart]
	sfg     singleflight.Group
	timeout time.Duration
}

type LRUCacheConfig struct {
	Size    int
	TTL     time.Duration
	Timeout time.Duration
}

func NewLRUCache(storage pcPartRepository, cfg *LRUCacheConfig) *LRUCache {
	const recentSize = 10
	return &LRUCache{
		storage: storage,
		recent:  expirable.NewLRU[string, []*domain.PcPart](recentSize, nil, cfg.TTL),
		inner:   expirable.NewLRU[string, *domain.PcPart](cfg.Size, nil, cfg.TTL),
		timeout: cfg.Timeout,
	}
}

func (c *LRUCache) CreatePcPart(ctx context.Context, part *domain.PcPart) error {
	return c.storage.CreatePcPart(ctx, part)
}

func (c *LRUCache) GetPcPartsRecent(ctx context.Context, limit int32) ([]*domain.PcPart, error) {
	const op = "cache.GetPcPartsRecent"
	key := fmt.Sprintf("recent:%d", limit)

	if cached, ok := c.recent.Get(key); ok {
		return cached, nil
	}

	data, err, _ := c.sfg.Do(key, func() (any, error) {
		if cached, ok := c.recent.Get(key); ok {
			return cached, nil
		}

		// Создаем контекст, который не умрет, если юзер нажал "отмена"
		asyncCtx := context.WithoutCancel(ctx)

		// Но обязательно добавляем свой таймаут, чтобы запрос не завис в БД навсегда
		asyncCtx, cancel := context.WithTimeout(asyncCtx, c.timeout)
		defer cancel()

		fetched, err := c.storage.GetPcPartsRecent(asyncCtx, limit)
		if err != nil {
			return nil, err
		}

		c.recent.Add(key, fetched)
		return fetched, nil
	})
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, ok := data.([]*domain.PcPart)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected type: %T", op, data)
	}

	return res, nil
}

func (c *LRUCache) GetPcPartByID(ctx context.Context, id uuid.UUID) (*domain.PcPart, error) {
	const op = "cache.GetPcPartByID"
	key := fmt.Sprintf("id:%s", id)

	if cached, ok := c.inner.Get(key); ok {
		return cached, nil
	}

	data, err, _ := c.sfg.Do(key, func() (any, error) {
		if cached, ok := c.inner.Get(key); ok {
			return cached, nil
		}

		asyncCtx := context.WithoutCancel(ctx)
		asyncCtx, cancel := context.WithTimeout(asyncCtx, c.timeout)
		defer cancel()

		fetched, err := c.storage.GetPcPartByID(asyncCtx, id)
		if err != nil {
			return nil, err
		}

		c.inner.Add(key, fetched)
		return fetched, nil
	})
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, ok := data.(*domain.PcPart)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected type: %T", op, data)
	}

	return res, nil
}

func (c *LRUCache) UpdatePcPart(ctx context.Context, part *domain.PcPart) error {
	return c.storage.UpdatePcPart(ctx, part)
}

func (c *LRUCache) SoftDeletePcPart(ctx context.Context, part *domain.PcPart) error {
	return c.storage.SoftDeletePcPart(ctx, part)
}
