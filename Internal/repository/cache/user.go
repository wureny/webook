package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wureny/webook/webook/Internal/domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (u *UserCache) Key(id uint64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (u *UserCache) Get(ctx context.Context, id uint64) (domain.User, error) {
	key := u.Key(id)
	val, err := u.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var s domain.User
	err = json.Unmarshal(val, &s)
	return s, err
}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.Key(u.Id)
	return cache.client.Set(ctx, key, string(val), cache.expiration).Err()
}
