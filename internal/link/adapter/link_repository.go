package adapter

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"reflect"
	"shortlink/internal/common/cache"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/error_no"
	"shortlink/internal/common/lock"
	"shortlink/internal/link/adapter/assembler"
	"shortlink/internal/link/adapter/po"
	"shortlink/internal/link/domain/link"
	"time"
)

type LinkRepository struct {
	db               *gorm.DB
	distributedCache cache.DistributedCache
	locker           lock.DistributedLock
	assembler        assembler.LinkAssembler
}

func NewLinkRepository(
	db *gorm.DB,
	distributedCache cache.DistributedCache,
	locker lock.DistributedLock,
) LinkRepository {
	return LinkRepository{
		db:               db,
		distributedCache: distributedCache,
		locker:           locker,
		assembler:        assembler.LinkAssembler{},
	}
}

func (r LinkRepository) ShortUriExists(ctx context.Context, shortUri string) (bool, error) {
	return r.distributedCache.ExistsInBloomFilter(
		ctx,
		cache.ShortUriCreateBloomFilter,
		constant.LockGotoLinkKey+shortUri,
		constant.GotoIsNullLinkKey+shortUri,
	)
}

//func (r LinkRepository) GetOriginalUrlByShortUrl(
//	ctx context.Context,
//	shortUri string,
//) (status int, res string, err error) {
//	// 尝试从缓存中获取短链的原始链接，存在则直接返回
//	value := r.rdb.Get(ctx, fmt.Sprintf(constant.GotoLinkKey, shortUri)).String()
//	if value != "" {
//		return value, nil
//	}
//
//	// 如果缓存中没有，判断是否存在于布隆过滤器中，不存在返回 not found
//	exists, err := r.rdb.BFExists(ctx, cache.ShortUriCreateBloomFilter, shortUri).Result()
//	if err != nil {
//		return "", err
//	}
//	if !exists {
//		// TODO 应用层需要处理这个错误，返回 404
//		return "", error_no.NewServiceError(error_no.LinkNotExists)
//	}
//
//	// 从缓存中获取 GotoIsNullLink 的值，如果存在则意味着短链接失效，返回 not found
//	gotoIsNullLink := r.rdb.Get(ctx, fmt.Sprintf(constant.GotoIsNullLinkKey, shortUri)).String()
//	if gotoIsNullLink != "" {
//		// TODO 应用层需要处理这个错误，返回 404
//		return "", error_no.NewServiceError(error_no.LinkNotExists)
//	}
//
//	// 获取分布式锁，并进行二次锁判定，尝试从数据库中获取原始链接，并写入缓存
//	lockKey := fmt.Sprintf(constant.LockGotoLinkKey, shortUri)
//	_, err = r.locker.Acquire(ctx, lockKey, time.Second)
//	if errors.Is(err, redislock.ErrNotObtained) {
//		return "", error_no.NewExternalErrorWithMsg(error_no.RedisError, "获取锁失败")
//	}
//	defer func() {
//		if err = r.locker.Release(ctx, lockKey); err != nil {
//			log.Errorf("释放锁失败: %v", err)
//		}
//	}()
//
//	// ------------- 在查询数据库之前再进行一次判断，防止缓存击穿 -------------
//	// 在高并发场景下，会存在多个线程同时竞争锁，但只有一个线程能够获取锁
//	// 第一个线程执行结束后，其他线程再次尝试获取锁，此时缓存中已经有值，直接返回
//	// TODO: 优化方案，可以使用分布式锁的续租功能，避免锁过期导致的缓存击穿
//	// TODO: 目前我是以Java的思维来写的，Go的锁机制可能有更好的解决方案
//	value = r.rdb.Get(ctx, fmt.Sprintf(constant.GotoLinkKey, shortUri)).String()
//	if value != "" {
//		return value, nil
//	}
//
//	gotoIsNullLink = r.rdb.Get(ctx, fmt.Sprintf(constant.GotoIsNullLinkKey, shortUri)).String()
//	if gotoIsNullLink != "" {
//		// TODO 应用层需要处理这个错误，返回 404
//		return "", error_no.NewServiceError(error_no.LinkNotExists)
//	}
//
//	// 查询数据库
//	var linkGoto po.LinkGoto
//	err = r.db.WithContext(ctx).
//		Model(&linkGoto).
//		Where("full_short_url = ? AND enable_status = 0 AND del_flag = false", shortUri).
//		First(&linkGoto).Error
//	if err != nil {
//		return "", err
//	}
//	// 数据库中不存在这个短链接
//	if linkGoto.FullShortUrl == "" {
//		r.rdb.SetEx(ctx, fmt.Sprintf(constant.GotoIsNullLinkKey, shortUri), "-", 30*time.Minute)
//		// TODO 应用层需要处理这个错误，返回 404
//		return "", error_no.NewServiceError(error_no.LinkNotExists)
//	}
//	// 数据库中存在短链接，则查询link表获取原始链接
//	var link po.Link
//	err = r.db.WithContext(ctx).
//		Model(&link).
//		Where("full_short_url = ? AND gid = ? AND enable_status = true AND del_flag = false", shortUri, linkGoto.Gid).
//		First(&link).Error
//	if err != nil {
//		return "", err
//	}
//	// 短链接失效
//	if link.FullShortUrl == "" || link.ValidDate.Before(time.Now()) {
//		// 写入缓存
//		r.rdb.SetEx(ctx, fmt.Sprintf(constant.GotoIsNullLinkKey, shortUri), "-", 30*time.Minute)
//		// TODO 应用层需要处理这个错误，返回 404
//		return "", error_no.NewServiceError(error_no.LinkNotExists)
//	}
//	// 查询到有效的原始链接，写入缓存
//	err = r.rdb.SetEx(ctx, fmt.Sprintf(constant.GotoLinkKey, shortUri), link.OriginUrl, toolkit.GetLinkCacheExpiration(link.ValidDate)).Err()
//	if err != nil {
//		return "", error_no.NewExternalErrorWithMsg(error_no.RedisError, err.Error())
//	}
//	return "", err
//}

// RecordLinkVisitInfo 记录短链接访问信息
//func (r LinkRepository) RecordLinkVisitInfo(ctx context.Context, info valobj.LinkStatssRecordVo) error {
//	// 确定两个值的信息，uvFirstFlag 和 uipFirstFlag
//	uvAdded, err := r.rdb.SAdd(ctx, consts.LinkStatsUvKey+info.ShortUri, info.UV).Result()
//	if err != nil {
//		return err
//	}
//	if uvAdded > 0 {
//		info.UVFirstFlag = true
//	}
//
//	uipAdded, err := r.rdb.SAdd(ctx, consts.LinkStatsUipKey+info.ShortUri, info.RemoteAddr).Result()
//	if err != nil {
//		return err
//	}
//	if uipAdded > 0 {
//		info.UipFirstFlag = true
//	}
//
//	msg := map[string]interface{}{
//		"statsRecord": info,
//	}
//	jsonMsg, err := sonic.Marshal(msg)
//	if err != nil {
//		return err
//	}
//	//_, err = r.producer.Send(ctx, &rmqclient.Message{
//	//	Topic: "app_short_link",
//	//	Tag:   nil,
//	//	Body:  jsonMsg,
//	//})
//	//if err != nil {
//	//	return err
//	//}
//	return nil
//}

// CreateLink 保存短链接并进行预热
func (r LinkRepository) CreateLink(ctx context.Context, lk *link.Link) (err error) {
	linkPo := r.assembler.LinkEntityToLinkPo(lk)
	linkGotoPo := r.assembler.LinkEntityToLinkGotoPo(lk)

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&linkPo).Error; err != nil {
			// 在高并发场景下可能出现重复插入的情况
			//if errors.Is(err, gorm.ErrDuplicatedKey) {
			//	// 添加到布隆过滤器
			//	if err = r.rdb.BFAdd(ctx, cache.ShortUriCreateBloomFilter, shortUri).Err(); err != nil {
			//		return error_no.NewServiceErrorWithMsg(error_no.LinkDuplicateInsert, err.Error())
			//	}
			//}
			return err
		}
		if err = tx.Create(&linkGotoPo).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}

	cacheValue := link.NewCacheValue(lk)
	err = r.distributedCache.SafePut(
		ctx,
		constant.GotoLinkKey+lk.ShortUri(),
		cacheValue,
		cacheValue.Expiration(),
		cache.ShortUriCreateBloomFilter,
		lk.ShortUri(),
	)
	return
}

func (r LinkRepository) CreateLinkBatch(ctx context.Context, links []*link.Link) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, lk := range links {
			linkPo := r.assembler.LinkEntityToLinkPo(lk)
			linkGotoPo := r.assembler.LinkEntityToLinkGotoPo(lk)

			if err := tx.Create(&linkPo).Error; err != nil {
				return err
			}
			if err := tx.Create(&linkGotoPo).Error; err != nil {
				return err
			}
		}

		for _, lk := range links {
			err := r.distributedCache.SafePut(
				ctx,
				constant.GotoLinkKey+lk.ShortUri(),
				link.NewCacheValue(lk),
				lk.ValidDate().Expiration(),
				cache.ShortUriCreateBloomFilter,
				lk.ShortUri(),
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateLink 更新短链接
// 1. 如果分组发生变化，需要删除原来的分组
// 2. 更新缓存
func (r LinkRepository) UpdateLink(
	ctx context.Context,
	id link.Identifier,
	updateFn func(ctx context.Context, link *link.Link) (*link.Link, error),
) (err error) {

	// 查询
	var linkPo po.Link
	if err = r.db.Where("gid = ? AND full_short_url = ?", id.Gid, id.ShortUri).
		First(&linkPo).Error; err != nil {
		return err
	}

	lk := r.assembler.LinkPoToLinkEntity(&linkPo)
	if lk, err = updateFn(ctx, lk); err != nil {
		return err
	}
	updatedLinkPo := r.assembler.LinkEntityToLinkPo(lk)

	if updatedLinkPo.Gid != linkPo.Gid {
		// 获取分布式锁
		lockKey := constant.LockGidUpdateKey + linkPo.ShortUri
		acquired := false
		if acquired, err = r.locker.Acquire(ctx, lockKey, constant.DefaultTimeOut); err != nil {
			return err
		}
		if !acquired {
			return error_no.LockAcquireFailed
		}

		// 事务更新
		err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// 从原来的分组中删除
			if err = tx.WithContext(ctx).Delete(&linkPo).Error; err != nil {
				return err
			}
			oldLinkGotoPo := po.LinkGoto{
				Gid:      linkPo.Gid,
				ShortUri: linkPo.ShortUri,
			}
			if err = tx.WithContext(ctx).Delete(oldLinkGotoPo).Error; err != nil {
				return err
			}

			// 插入新的分组
			if err = tx.WithContext(ctx).Create(&updatedLinkPo).Error; err != nil {
				return err
			}
			shortLinkGotoPo := po.LinkGoto{
				Gid:      updatedLinkPo.Gid,
				ShortUri: updatedLinkPo.ShortUri,
			}
			if err = tx.WithContext(ctx).Create(&shortLinkGotoPo).Error; err != nil {
				return err
			}
			return nil
		})

		if err = r.locker.Release(ctx, lockKey); err != nil {
			log.Errorf("释放锁失败: %v", err)
			return err
		}
	} else {
		if err = r.db.Save(&updatedLinkPo).Error; err != nil {
			return err
		}
	}

	// 更新缓存
	err = r.distributedCache.Put(
		ctx,
		constant.GotoLinkKey+lk.ShortUri(),
		link.NewCacheValue(lk),
		lk.ValidDate().Expiration(),
	)
	return err
}

// SaveToRecycleBin 保存到回收站
func (r LinkRepository) SaveToRecycleBin(
	ctx context.Context,
	id link.Identifier,
) (err error) {

	// 持久化操作
	var linkPo po.Link
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询
		if err = r.db.WithContext(ctx).
			Model(&linkPo).Where("gid = ? and shortUri = ?", id.Gid, id.ShortUri).
			First(&linkPo).Error; err != nil {
			return err
		}

		// 修改
		linkPo.RecycleTime = sql.NullTime{
			Time: time.Now(),
		}

		// 更新
		if err = r.db.WithContext(ctx).Save(&linkPo).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_no.LinkNotExists
		}
		return err
	}

	// 修改缓存中的状态为已删除
	cacheKey := constant.GotoLinkKey + id.ShortUri
	if err = r.modifyCacheValueStatus(ctx, cacheKey, link.StatusDeleted); err != nil {
		return
	}
	return
}

// RemoveFromRecycleBin 从回收站移除
func (r LinkRepository) RemoveFromRecycleBin(
	ctx context.Context,
	id link.Identifier,
) (err error) {

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		var linkPo po.Link
		err = r.db.WithContext(ctx).Model(&linkPo).
			Where("gid = ? and fullShortUrl = ?", id.Gid, id.ShortUri).
			Find(&linkPo).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return error_no.LinkNotExists
			}
			return err
		}

		//如果短链接状态不是回收站状态，返回错误
		if !linkPo.RecycleTime.Valid {
			return error_no.InvalidLinkStatsus
		}

		err = r.db.Delete(&linkPo).Error
		return
	})

	// 删除缓存
	err = r.distributedCache.SafeDelete(
		ctx, constant.GotoLinkKey+id.ShortUri,
		constant.GotoIsNullLinkKey+id.ShortUri,
	)
	if err != nil {
		return err
	}

	return
}

// RecoverFromRecycleBin 从回收站中恢复
func (r LinkRepository) RecoverFromRecycleBin(
	ctx context.Context,
	id link.Identifier,
) (err error) {

	var linkPo po.Link
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		err = r.db.WithContext(ctx).Model(&linkPo).
			Where("gid = ? and fullShortUrl = ?", id.Gid, id.ShortUri).
			Find(&linkPo).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return error_no.LinkNotExists
			}
			return
		}

		// 如果短链接状态不是回收站状态，返回错误
		if !linkPo.RecycleTime.Valid {
			return error_no.InvalidLinkStatsus
		}

		// 修改
		linkPo.RecycleTime = sql.NullTime{}

		err = r.db.Save(&linkPo).Error
		return
	})
	if err != nil {
		return
	}

	cacheKey := constant.GotoLinkKey + id.ShortUri
	if err = r.modifyCacheValueStatus(ctx, cacheKey, linkPo.Status); err != nil {
		return
	}
	return
}

func (r LinkRepository) modifyCacheValueStatus(ctx context.Context, cacheKey string, status string) (err error) {
	var res interface{}
	cacheValue := link.CacheValue{}
	if res, err = r.distributedCache.Get(ctx, cacheKey, reflect.TypeOf(link.CacheValue{})); err != nil {
		return err
	}
	cacheValue = res.(link.CacheValue)
	cacheValue.Status = status
	if err = r.distributedCache.Put(ctx, cacheKey, cacheValue, cacheValue.Expiration()); err != nil {
		return err
	}
	return
}
