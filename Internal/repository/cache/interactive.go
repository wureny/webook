package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/wureny/webook/webook/Internal/domain"
	"time"
)

type InteractiveCache interface {

	// IncrReadCntIfPresent 如果在缓存中有对应的数据，就 +1
	IncrReadCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	// Get 查询缓存中数据
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
}

type RedisInteractiveCache struct {
	client redis.Cmdable
	// 缓存过期时间
	expiration time.Duration
}

func NewRedisInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		client: client,
	}
}

func (r RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisInteractiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisInteractiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisInteractiveCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
	//TODO implement me
	panic("implement me")
}
