package repo

import (
	"context"
	"errors"
	"fmt"
	rmqclient "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/bsm/redislock"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"shortlink/common/consts"
	"shortlink/common/error_no"
	"shortlink/common/toolkit"
	"shortlink/common/types"
	"shortlink/internal/domain/link/aggregate"
	"shortlink/internal/domain/link/valobj"
	"shortlink/internal/infra/cache"
	"shortlink/internal/infra/persistence/converter"
	po2 "shortlink/internal/infra/persistence/po"
	"time"
)

type LinkRepository struct {
	db       *gorm.DB
	rdb      *redis.Client
	locker   *redislock.Client
	producer rmqclient.Producer
	cvt      converter.LinkConverter
}

func NewLinkRepository(db *gorm.DB) LinkRepository {
	return LinkRepository{db: db}
}

func (r LinkRepository) GetOriginalUrlByShortUrl(
	ctx context.Context,
	shortUrl string,
	statsInfo valobj.ShortLinkStatsRecordVO,
) (res string, err error) {
	// 尝试从缓存中获取短链的原始链接，存在则直接返回
	value := r.rdb.Get(ctx, fmt.Sprintf(consts.GotoShortLinkKey, shortUrl)).String()
	if value != "" {
		// 记录用户访问信息
		if err := r.recordStatsAndSetUser(ctx, statsInfo); err != nil {
			return "", err
		}
		return value, nil
	}

	// 如果缓存中没有，判断是否存在于布隆过滤器中，不存在返回 not found
	exists, err := r.rdb.BFExists(ctx, cache.ShortUriCreateBloomFilter, shortUrl).Result()
	if err != nil {
		return "", err
	}
	if !exists {
		// TODO 应用层需要处理这个错误，返回 404
		return "", error_no.NewServiceError(error_no.ShortLinkNotExists)
	}

	// 从缓存中获取 GotoIsNullShortLink 的值，如果存在则意味着短链接失效，返回 not found
	gotoIsNullShortLink := r.rdb.Get(ctx, fmt.Sprintf(consts.GotoIsNullShortLinkKey, shortUrl)).String()
	if gotoIsNullShortLink != "" {
		// TODO 应用层需要处理这个错误，返回 404
		return "", error_no.NewServiceError(error_no.ShortLinkNotExists)
	}

	// 获取分布式锁，并进行二次锁判定，尝试从数据库中获取原始链接，并写入缓存
	lock, err := r.locker.Obtain(ctx, fmt.Sprintf(consts.LockGotoShortLinkKey, shortUrl), time.Second, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return "", error_no.NewExternalErrorWithMsg(error_no.RedisError, "获取锁失败")
	}
	defer func(lock *redislock.Lock, ctx context.Context) {
		err := lock.Release(ctx)
		if err != nil {
			log.Errorf("释放锁失败: %v", err)
		}
	}(lock, ctx)

	// ------------- 在查询数据库之前再进行一次判断，防止缓存击穿 -------------
	// 在高并发场景下，会存在多个线程同时竞争锁，但只有一个线程能够获取锁
	// 第一个线程执行结束后，其他线程再次尝试获取锁，此时缓存中已经有值，直接返回
	// TODO: 优化方案，可以使用分布式锁的续租功能，避免锁过期导致的缓存击穿
	// TODO: 目前我是以Java的思维来写的，Go的锁机制可能有更好的解决方案
	value = r.rdb.Get(ctx, fmt.Sprintf(consts.GotoShortLinkKey, shortUrl)).String()
	if value != "" {
		// 记录用户访问信息
		if err := r.recordStatsAndSetUser(ctx, statsInfo); err != nil {
			return "", err
		}
		return value, nil
	}

	gotoIsNullShortLink = r.rdb.Get(ctx, fmt.Sprintf(consts.GotoIsNullShortLinkKey, shortUrl)).String()
	if gotoIsNullShortLink != "" {
		// TODO 应用层需要处理这个错误，返回 404
		return "", error_no.NewServiceError(error_no.ShortLinkNotExists)
	}

	// 查询数据库
	var linkGoto po2.LinkGoto
	err = r.db.WithContext(ctx).
		Model(&linkGoto).
		Where("full_short_url = ? AND enable_status = 0 AND del_flag = false", shortUrl).
		First(&linkGoto).Error
	if err != nil {
		return "", err
	}
	// 数据库中不存在这个短链接
	if linkGoto.FullShortURL == "" {
		r.rdb.SetEx(ctx, fmt.Sprintf(consts.GotoIsNullShortLinkKey, shortUrl), "-", 30*time.Minute)
		// TODO 应用层需要处理这个错误，返回 404
		return "", error_no.NewServiceError(error_no.ShortLinkNotExists)
	}
	// 数据库中存在短链接，则查询link表获取原始链接
	var link po2.Link
	err = r.db.WithContext(ctx).
		Model(&link).
		Where("full_short_url = ? AND gid = ? AND enable_status = true AND del_flag = false", shortUrl, linkGoto.Gid).
		First(&link).Error
	if err != nil {
		return "", err
	}
	// 短链接失效
	if link.FullShortUrl == "" || link.ValidDate.Before(time.Now()) {
		// 写入缓存
		r.rdb.SetEx(ctx, fmt.Sprintf(consts.GotoIsNullShortLinkKey, shortUrl), "-", 30*time.Minute)
		// TODO 应用层需要处理这个错误，返回 404
		return "", error_no.NewServiceError(error_no.ShortLinkNotExists)
	}
	// 查询到有效的原始链接，写入缓存
	err = r.rdb.SetEx(ctx, fmt.Sprintf(consts.GotoShortLinkKey, shortUrl), link.OriginalUrl, toolkit.GetLinkCacheExpiration(link.ValidDate)).Err()
	if err != nil {
		return "", error_no.NewExternalErrorWithMsg(error_no.RedisError, err.Error())
	}
	// 记录用户访问信息
	if err := r.recordStatsAndSetUser(ctx, statsInfo); err != nil {
		return "", err
	}
	return "", err
}

func (r LinkRepository) recordStatsAndSetUser(ctx context.Context, statsInfo valobj.ShortLinkStatsRecordVO) error {
	message := map[string]interface{}{
		"statsRecord": statsInfo,
	}
	_, err := r.producer.Send(ctx, &rmqclient.Message{
		Topic: "shortlink",
		Tag:   nil,
		Body:  []byte(fmt.Sprintf("%v", message)),
	})
	if err != nil {
		return error_no.NewExternalErrorWithMsg(error_no.RocketMQError, err.Error())
	}
	return nil
}

// CreateLink 保存短链接并进行预热
func (r LinkRepository) CreateLink(ctx context.Context, aggregate aggregate.CreateLinkAggregate) (err error) {
	shortLink := aggregate.ShortLink
	shortLinkGoTo := aggregate.ShortLinkGoto
	shortUrl := shortLink.ShortUrl()

	// 生成短链接
	//shortUri, err := r.generateUniqueShortLink(ctx, shortLink.GenerateFunc(), 10)
	//if err != nil {
	//	return err
	//}
	//defaultDomain := config.ShortLinkDomain.String()
	//fullShortUrl := defaultDomain + "/" + shortUri

	// 领域实体赋值
	//shortLink.Domain = defaultDomain
	//shortLink.ShortUri = shortUri
	//shortLink.FullShortUrl = fullShortUrl
	//shortLinkGoTo.FullShortUrl = fullShortUrl

	// 持久化短链接
	err = r.db.Transaction(func(tx *gorm.DB) error {
		shortLinkPo := r.cvt.LinkEntityToPo(&shortLink)
		shortLinkGotoPo := r.cvt.LinkGotoEntityToPo(shortLinkGoTo)

		if err := tx.Create(&shortLinkPo).Error; err != nil {
			// 在高并发场景下可能出现重复插入的情况
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				// 添加到布隆过滤器
				if err := r.rdb.BFAdd(ctx, cache.ShortUriCreateBloomFilter, shortUrl).Err(); err != nil {
					return error_no.NewServiceErrorWithMsg(error_no.ShortLinkDuplicateInsert, err.Error())
				}
			}
			return err
		}
		if err := tx.Create(&shortLinkGotoPo).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 预热短链接
	key := fmt.Sprintf(consts.GotoShortLinkKey, shortUrl)
	expiration := toolkit.GetLinkCacheExpiration(shortLink.ValidDate())
	err = r.rdb.SetEx(ctx, key, shortUrl, expiration).Err()
	if err != nil {
		return error_no.NewExternalErrorWithMsg(error_no.RedisError, err.Error())
	}

	// 添加到布隆过滤器
	err = r.rdb.BFAdd(ctx, cache.ShortUriCreateBloomFilter, shortUrl).Err()
	if err != nil {
		return error_no.NewExternalErrorWithMsg(error_no.RedisError, err.Error())
	}
	return nil
}

// CreateLinkWithLock 通过分布式锁创建短链接
func (r LinkRepository) CreateLinkWithLock(ctx context.Context, aggregate aggregate.CreateLinkAggregate) error {
	// 获取分布式锁
	lock, err := r.locker.Obtain(ctx, consts.ShortLinkCreateLockKey, time.Second, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return error_no.NewExternalErrorWithMsg(error_no.RedisError, "获取锁失败")
	} else if err != nil {
		return error_no.NewExternalErrorWithMsg(error_no.RedisError, err.Error())
	}
	// 释放锁
	defer func(lock *redislock.Lock, ctx context.Context) {
		err := lock.Release(ctx)
		if err != nil {
			log.Errorf("释放锁失败: %v", err)
		}
	}(lock, ctx)

	return r.CreateLink(ctx, aggregate)
}

// UpdateLink 更新短链接
// 1. 如果分组发生变化，需要删除原来的分组
// 2. 更新缓存
func (r LinkRepository) UpdateLink(
	ctx context.Context,
	id types.LinkID,
	enableStatus int,
	updateFn func(ctx context.Context, link *types.Link) (*types.Link, error),
) (err error) {

	// 查询
	var linkPo po2.Link
	if err = r.db.Where("gid = ? AND full_short_url = ?", id.Gid, id.FullShortUrl).
		Where("enable_status = ?", enableStatus).
		First(&linkPo).Error; err != nil {
		return err
	}

	link := r.cvt.LinkPoToEntity(linkPo)
	link, err = updateFn(ctx, link)
	if err != nil {
		return err
	}
	updatedLinkPo := r.cvt.LinkEntityToPo(link)

	if updatedLinkPo.Gid != linkPo.Gid {
		// 获取分布式锁
		lockKey := fmt.Sprintf(consts.LockGidUpdateKey, linkPo.FullShortUrl)
		lock, err := r.locker.Obtain(ctx, lockKey, 100*time.Millisecond, nil)
		if errors.Is(err, redislock.ErrNotObtained) {
			return error_no.NewExternalErrorWithMsg(error_no.RedisError, "获取锁失败")
		}

		// 事务更新
		err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// 从原来的分组中删除
			err = r.db.WithContext(ctx).Delete(&linkPo).Error
			if err != nil {
				return err
			}
			oldShortLinkGotoPo := po2.LinkGoto{
				Gid:          linkPo.Gid,
				FullShortURL: linkPo.FullShortUrl,
			}
			err = r.db.WithContext(ctx).Delete(oldShortLinkGotoPo).Error
			if err != nil {
				return err
			}

			// 插入新的分组
			if err := tx.Create(&updatedLinkPo).Error; err != nil {
				return err
			}
			shortLinkGotoPo := po2.LinkGoto{
				Gid:          updatedLinkPo.Gid,
				FullShortURL: updatedLinkPo.FullShortUrl,
			}
			if err := tx.Create(&shortLinkGotoPo).Error; err != nil {
				return err
			}
			return nil
		})

		err = lock.Release(ctx)
		if err != nil {
			log.Errorf("释放锁失败: %v", err)
		}
	} else {
		err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err = tx.Save(&updatedLinkPo).Error; err != nil {
				return err
			}
			return nil
		})
	}

	if err := r.db.Save(&updatedLinkPo).Error; err != nil {
		return err
	}

	// 更新缓存
	if linkPo.ValidDateType != updatedLinkPo.ValidDateType ||
		linkPo.ValidDate != updatedLinkPo.ValidDate ||
		linkPo.OriginalUrl != updatedLinkPo.OriginalUrl {
		key := fmt.Sprintf(consts.GotoShortLinkKey, linkPo.FullShortUrl)
		err = r.rdb.Del(ctx, key).Err()
		if err != nil {
			return error_no.NewExternalErrorWithMsg(error_no.RedisError, err.Error())
		}
	}
	if linkPo.ValidDate.Before(time.Now()) && (updatedLinkPo.ValidDate.After(time.Now()) || updatedLinkPo.ValidDateType == 0) {
		key := fmt.Sprintf(consts.GotoIsNullShortLinkKey, linkPo.FullShortUrl)
		r.rdb.Del(ctx, key)
	}

	return nil
}

func (r LinkRepository) generateUniqueShortLink(ctx context.Context, gen func() string, attempts int) (string, error) {
	count, result := 0, ""
	for {
		if count > attempts {
			return "", error_no.NewServiceError(error_no.TooManyShortLinkCreate)
		}
		result = gen()
		exists, err := r.rdb.BFExists(ctx, cache.ShortUriCreateBloomFilter, result).Result()
		if err != nil {
			return "", error_no.NewExternalErrorWithMsg(error_no.RedisError, err.Error())
		}
		if !exists {
			break
		}
		count++
	}
	return result, nil
}