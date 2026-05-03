package pcpart

import (
	"context"
	"fmt"
	"time"

	"github.com/KlementevTech/gotips/internal/domain"
	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/valyala/fastrand"
	"golang.org/x/sync/singleflight"
)

type shard struct {
	recent *expirable.LRU[string, []*domain.PcPart]
	items  *expirable.LRU[string, *domain.PcPart]
}

type LRUCache struct {
	shardsCount uint32
	shards      []*shard
	sfg         singleflight.Group
	storage     domain.Repository
	timeout     time.Duration
}

type LRUCacheConfig struct {
	Size    int
	TTL     time.Duration
	Timeout time.Duration
	Shards  uint32
}

func NewLRUCache(storage domain.Repository, cfg *LRUCacheConfig) *LRUCache {
	const recentsCnt = 5

	shards := make([]*shard, cfg.Shards)
	for i := range shards {
		shards[i] = &shard{
			recent: expirable.NewLRU[string, []*domain.PcPart](recentsCnt, nil, cfg.TTL),
			items:  expirable.NewLRU[string, *domain.PcPart](cfg.Size, nil, cfg.TTL),
		}
	}

	return &LRUCache{
		storage:     storage,
		shardsCount: cfg.Shards,
		shards:      shards,
		timeout:     cfg.Timeout,
	}
}

func (c *LRUCache) CreatePcPart(ctx context.Context, params domain.CreatePcPartParams) (*domain.PcPart, error) {
	return c.storage.CreatePcPart(ctx, params)
}

func (c *LRUCache) GetPcPartsRecent(ctx context.Context, limit int32) ([]*domain.PcPart, error) {
	const op = "cache.GetPcPartsRecent"
	key := fmt.Sprintf("recent:%d", limit)

	shrd := c.shards[fastrand.Uint32n(c.shardsCount)]

	if cached, ok := shrd.recent.Get(key); ok {
		return cached, nil
	}

	data, err, _ := c.sfg.Do(key, func() (any, error) {
		if cached, ok := shrd.recent.Get(key); ok {
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

		shrd.recent.Add(key, fetched)
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
	key := fmt.Sprintf("id:%d", id.ID())

	// Если c.shardsCount - степень 2, то % можно заменить на сдвиг
	shrd := c.shards[id.ID()&(c.shardsCount-1)]

	if cached, ok := shrd.items.Get(key); ok {
		return cached, nil
	}

	data, err, _ := c.sfg.Do(key, func() (any, error) {
		if cached, ok := shrd.items.Get(key); ok {
			return cached, nil
		}

		asyncCtx := context.WithoutCancel(ctx)
		asyncCtx, cancel := context.WithTimeout(asyncCtx, c.timeout)
		defer cancel()

		fetched, err := c.storage.GetPcPartByID(asyncCtx, id)
		if err != nil {
			return nil, err
		}

		shrd.items.Add(key, fetched)
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

func (c *LRUCache) UpdatePcPart(ctx context.Context, params domain.UpdatePcPartParams) (*domain.PcPart, error) {
	return c.storage.UpdatePcPart(ctx, params)
}

func (c *LRUCache) SoftDeletePcPart(ctx context.Context, id uuid.UUID, version int) (*domain.PcPart, error) {
	return c.storage.SoftDeletePcPart(ctx, id, version)
}
