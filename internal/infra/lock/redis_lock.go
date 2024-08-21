package lock

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisLock struct {
	locker *redislock.Client
}

func NewRedisLock(rdb *redis.Client) RedisLock {
	return RedisLock{locker: redislock.New(rdb)}
}

var Locks = map[string]*redislock.Lock{}

func (r RedisLock) Acquire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	obtain, err := r.locker.Obtain(ctx, key, expiration, nil)
	if err != nil {
		return false, err
	}
	Locks[key] = obtain
	return true, nil
}

func (r RedisLock) Release(ctx context.Context, key string) error {
	if err := Locks[key].Release(ctx); err != nil {
		return err
	}
	delete(Locks, key)
	return nil
}

func (r RedisLock) Refresh(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	if err := Locks[key].Refresh(ctx, expiration, nil); err != nil {
		return false, err
	}
	return true, nil
}
