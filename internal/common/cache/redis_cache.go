package cache

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"reflect"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/error_no"
	"shortlink/internal/common/lock"
	"time"
)

type RedisDistributedCache struct {
	rdb    *redis.Client
	locker lock.DistributedLock
}

func NewRedisDistributedCache(rdb *redis.Client, locker lock.DistributedLock) *RedisDistributedCache {
	if rdb == nil {
		panic("nil rdb")
	}
	if locker == nil {
		panic("nil locker")
	}
	return &RedisDistributedCache{rdb: rdb, locker: locker}
}

func (r RedisDistributedCache) Get(ctx context.Context, key string, valueType reflect.Type) (interface{}, error) {
	value, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, error_no.RedisKeyNotExist
		}
		return nil, err
	}

	if valueType == reflect.TypeOf("") {
		return value, nil
	}

	result := reflect.New(valueType).Interface()
	if err = sonic.Unmarshal([]byte(value), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r RedisDistributedCache) Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	valueBytes, err := sonic.Marshal(value)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, key, string(valueBytes), expiration).Err()
}

func (r RedisDistributedCache) PutIfAbsent(ctx context.Context, key string, value interface{}) (bool, error) {
	luaScript := `
		if redis.call("EXISTS", KEYS[1]) == 0 then
			redis.call("SET", KEYS[1], ARGV[1])
			return 1
		else
			return 0
		end
	`

	valueBytes, err := sonic.Marshal(value)
	if err != nil {
		return false, err
	}

	result, err := r.rdb.Eval(ctx, luaScript, []string{key}, valueBytes).Result()
	if err != nil {
		return false, err
	}

	if result.(int) == 1 {
		return true, nil
	}
	return false, nil
}

func (r RedisDistributedCache) Delete(ctx context.Context, key string) (bool, error) {
	result, err := r.rdb.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r RedisDistributedCache) DeleteMultiple(ctx context.Context, keys []string) (int, error) {
	result, err := r.rdb.Del(ctx, keys...).Result()
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

func (r RedisDistributedCache) HasKey(ctx context.Context, key string) (bool, error) {
	result, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
func (r RedisDistributedCache) GetInstance() interface{} {
	return r.rdb
}

func (r RedisDistributedCache) SafeGet(
	ctx context.Context,
	key string,
	valueType reflect.Type,
	cacheLoader Loader,
	expiration time.Duration,
) (interface{}, error) {
	return r.SafeGetWithCacheGetIfAbsent(ctx, key, valueType, cacheLoader, expiration, "", "", "", nil)
}

func (r RedisDistributedCache) SafeGetWithBloomFilter(
	ctx context.Context,
	key string,
	valueType reflect.Type,
	cacheLoader Loader,
	expiration time.Duration,
	bloomFilter string,
	bloomKey string,
) (interface{}, error) {
	return r.SafeGetWithCacheGetIfAbsent(ctx, key, valueType, cacheLoader, expiration, bloomFilter, bloomKey, "", nil)
}

func (r RedisDistributedCache) SafeGetWithCacheCheckFilter(
	ctx context.Context,
	key string,
	valueType reflect.Type,
	cacheLoader Loader,
	expiration time.Duration,
	bloomFilter string,
	bloomKey string,
	exceptBloomKey string,
) (interface{}, error) {
	return r.SafeGetWithCacheGetIfAbsent(ctx, key, valueType, cacheLoader, expiration, bloomFilter, bloomKey, exceptBloomKey, nil)
}

func (r RedisDistributedCache) SafeGetWithCacheGetIfAbsent(
	ctx context.Context,
	key string,
	valueType reflect.Type,
	cacheLoader Loader,
	expiration time.Duration,
	bloomFilter string,
	bloomKey string,
	exceptBloomKey string,
	cacheGetIfAbsent GetIfAbsent,
) (interface{}, error) {
	// step1 从缓存中取值
	result, err := r.Get(ctx, key, valueType)
	if err != nil {
		if errors.Is(err, error_no.RedisKeyNotExist) {

		} else {
			return nil, err
		}
	}
	_isNilOrEmpty := false
	if _isNilOrEmpty, err = isNilOrEmpty(result); err != nil {
		return nil, err
	}
	if !_isNilOrEmpty {
		return result, nil
	}

	// 缓存中为空 可能是因为缓存失效 判断是否存在于布隆过滤器中
	inBloom, deleteFromBloom := true, false
	if bloomFilter != "" {
		if inBloom, err = rdb.BFExists(ctx, bloomFilter, bloomKey).Result(); err != nil {
			return nil, err
		}
	}
	if exceptBloomKey != "" {
		if _, err = rdb.Get(ctx, exceptBloomKey).Result(); err != nil {
			if errors.Is(err, redis.Nil) {
				deleteFromBloom = false
			} else {
				return nil, err
			}
		} else {
			deleteFromBloom = true
		}
	}
	// 存在于缓存 -> 返回值
	// 不在布隆过滤器 | 已经失效 -> 无需查询数据库
	if !inBloom || deleteFromBloom {
		return nil, error_no.RedisKeyNotExist
	}

	// step2 获取分布式锁
	acquired := false
	lockKey := constant.LockGotoLinkKey + key
	if acquired, err = r.locker.Acquire(ctx, lockKey, constant.DefaultTimeOut); err != nil {
		return result, err
	}
	if !acquired {
		return result, error_no.LockAcquireFailed
	}
	defer func(locker lock.DistributedLock, ctx context.Context, key string) {
		if releaseErr := locker.Release(ctx, key); releaseErr != nil {
			err = releaseErr
		}
	}(r.locker, ctx, lockKey)

	// 双重判断，防止缓存击穿
	if result, err = r.Get(ctx, key, valueType); err != nil {
		if errors.Is(err, error_no.RedisKeyNotExist) {

		} else {
			return nil, err
		}
	}
	if _isNilOrEmpty, err = isNilOrEmpty(result); err != nil {
		return nil, err
	}
	if _isNilOrEmpty {
		// 从数据库中获取
		if result, err = r.loadAndSet(ctx, key, cacheLoader, expiration, bloomFilter, bloomKey, false); err != nil {
			return nil, err
		}
		if _isNilOrEmpty, err = isNilOrEmpty(result); err != nil {
			return nil, err
		}
		if _isNilOrEmpty {
			if cacheGetIfAbsent != nil {
				if err = cacheGetIfAbsent(key); err != nil {
					return nil, err
				}
			}
		}
	}
	return result, nil
}

func (r RedisDistributedCache) SafePut(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
	bloomFilter string,
	bloomKey string,
) error {
	if err := r.Put(ctx, key, value, expiration); err != nil {
		return err
	}
	return r.rdb.BFAdd(ctx, bloomFilter, bloomKey).Err()
}

func (r RedisDistributedCache) SafeDelete(ctx context.Context, key string, exceptBloomKey string) error {
	if ok, err := r.Delete(ctx, key); err != nil {
		return err
	} else if !ok {
		return nil
	}
	if exceptBloomKey != "" {
		if err := r.Put(ctx, exceptBloomKey, "-", constant.NeverExpire); err != nil {
			return err
		}
	}
	return nil
}

func (r RedisDistributedCache) ExistsInBloomFilter(ctx context.Context, key string, bloomFilter string, exceptKey string) (bool, error) {
	if exceptKey != "" {
		if _, err := r.rdb.Get(ctx, exceptKey).Result(); err != nil {
			if errors.Is(err, redis.Nil) {

			} else {
				return false, err
			}
		}
		// 在失效缓存中，意味着从布隆过滤器中删除了
		return false, nil
	}
	return r.rdb.BFExists(ctx, bloomFilter, key).Result()
}
func (r RedisDistributedCache) CountExistingKeys(ctx context.Context, keys ...string) (int, error) {
	result, err := r.rdb.Exists(ctx, keys...).Result()
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

func (r RedisDistributedCache) loadAndSet(
	ctx context.Context,
	key string,
	cacheLoader Loader,
	expiration time.Duration,
	bloomFilter string,
	bloomKey string,
	safeFlag bool,
) (interface{}, error) {
	result, err := cacheLoader()
	if err != nil {
		return nil, err
	}
	_isNilOrEmpty := false
	if _isNilOrEmpty, err = isNilOrEmpty(result); err != nil {
		return result, err
	}
	if _isNilOrEmpty {
		return nil, err
	}
	if safeFlag {
		if err = r.SafePut(ctx, key, result, expiration, bloomFilter, bloomKey); err != nil {
			return nil, err
		}
	} else {
		if err = r.Put(ctx, key, result, expiration); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (r RedisDistributedCache) DoubleDelete(ctx context.Context, key string, delay time.Duration) error {
	return nil
}
