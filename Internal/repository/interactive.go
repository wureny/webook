package repository

import (
	"context"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/repository/cache"
	"github.com/wureny/webook/webook/Internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context,
		biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type CachedReadCntRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
}

func NewCachedReadCntRepository(cache cache.InteractiveCache, dao dao.InteractiveDAO) InteractiveRepository {
	return &CachedReadCntRepository{
		cache: cache,
		dao:   dao,
	}
}

func (c CachedReadCntRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	// 要考虑缓存方案了
	// 这两个操作能不能换顺序？ —— 不能
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	go func() {
		_ = c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
	}()
	return nil
}

func (c CachedReadCntRepository) IncrLike(ctx context.Context, biz string, bizId, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (c CachedReadCntRepository) DecrLike(ctx context.Context, biz string, bizId, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, bizId)

}

func (c CachedReadCntRepository) AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (c CachedReadCntRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	var res domain.Interactive
	var ress dao.Interactive
	ress, err := c.dao.Get(ctx, biz, bizId)
	if err != nil {
		return res, nil
	}
	res = c.toDomain(ress)
	return res, nil
}

func (c CachedReadCntRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c CachedReadCntRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedReadCntRepository) toDomain(intr dao.Interactive) domain.Interactive {
	return domain.Interactive{
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		ReadCnt:    intr.ReadCnt,
	}
}
