package readrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"shortlink/common/consts"
	"shortlink/common/error_no"
	"shortlink/common/toolkit"
	"shortlink/common/types"
	"shortlink/internal/app/link/query"
	"shortlink/internal/infra/cache"
	"shortlink/internal/infra/lock"
	"shortlink/internal/infra/persistence/readrepo/internal/dao"
	"time"
)

type LinkQuery struct {
	linkDao dao.LinkDao
	rdb     *redis.Client
	locker  lock.DistributedLock
}

func NewLinkQuery(db *gorm.DB) LinkQuery {
	return LinkQuery{linkDao: dao.NewLinkDao(db)}
}

// GetOriginalUrlByShortUrl 根据短链接获取原始链接
func (q LinkQuery) GetOriginalUrlByShortUrl(ctx context.Context, fullShortUrl string) (string, error) {
	// 尝试从缓存中获取短链的原始链接，存在则直接返回
	key := fmt.Sprintf(consts.GotoShortLinkKey, fullShortUrl)
	value := q.rdb.Get(ctx, key).String()
	if value != "" {
		return value, nil
	}

	// 如果缓存中没有，判断是否存在于布隆过滤器中，不存在返回 not found
	exists, err := q.rdb.BFExists(ctx, cache.ShortUriCreateBloomFilter, fullShortUrl).Result()
	if err != nil {
		return "", err
	}
	if !exists {
		return "", error_no.ShortLinkNotFound
	}

	// 从缓存中获取 GotoIsNullShortLink 的值，如果存在则意味着短链接失效，返回 not found
	key = fmt.Sprintf(consts.GotoIsNullShortLinkKey, fullShortUrl)
	gotoIsNullShortLink := q.rdb.Get(ctx, key).String()
	if gotoIsNullShortLink != "" {
		return "", error_no.ShortLinkExpired
	}

	// 获取分布式锁，并进行二次锁判定，尝试从数据库中获取原始链接，并写入缓存
	lockKey := fmt.Sprintf(consts.LockGotoShortLinkKey, fullShortUrl)
	_, err = q.locker.Acquire(ctx, lockKey, time.Second)
	if errors.Is(err, redislock.ErrNotObtained) {
		return "", error_no.LockAcquireFailed
	}
	defer func() {
		if err = q.locker.Release(ctx, lockKey); err != nil {
			err = error_no.LockReleaseFailed
		}
	}()

	// ------------- 在查询数据库之前再进行一次判断，防止缓存击穿 -------------
	// 在高并发场景下，会存在多个线程同时竞争锁，但只有一个线程能够获取锁
	// 第一个线程执行结束后，其他线程再次尝试获取锁，此时缓存中已经有值，直接返回
	// TODO: 优化方案，可以使用分布式锁的续租功能，避免锁过期导致的缓存击穿
	// TODO: 目前我是以Java的思维来写的，Go的锁机制可能有更好的解决方案
	key = fmt.Sprintf(consts.GotoShortLinkKey, fullShortUrl)
	value = q.rdb.Get(ctx, key).String()
	if value != "" {
		return value, nil
	}

	key = fmt.Sprintf(consts.GotoIsNullShortLinkKey, fullShortUrl)
	gotoIsNullShortLink = q.rdb.Get(ctx, key).String()
	if gotoIsNullShortLink != "" {
		return "", error_no.ShortLinkExpired
	}

	// 查询数据库
	linkGotoPo, err := q.linkDao.GetLinkGoto(ctx, fullShortUrl)
	if err != nil {
		return "", err
	}
	// 数据库中不存在这个短链接
	if linkGotoPo.FullShortURL == "" {
		key = fmt.Sprintf(consts.GotoIsNullShortLinkKey, fullShortUrl)
		q.rdb.SetEx(ctx, key, "-", 30*time.Minute)
		return "", error_no.ShortLinkNotFound
	}
	// 数据库中存在短链接，则查询link表获取原始链接
	linkPo, err := q.linkDao.GetLink(ctx, types.LinkID{FullShortUrl: fullShortUrl, Gid: linkGotoPo.Gid})
	if err != nil {
		return "", err
	}
	// 短链接失效
	if linkPo.FullShortUrl == "" || linkPo.ValidDate.Before(time.Now()) {
		// 写入缓存
		key = fmt.Sprintf(consts.GotoIsNullShortLinkKey, fullShortUrl)
		q.rdb.SetEx(ctx, key, "-", 30*time.Minute)
		return "", error_no.ShortLinkExpired
	}
	// 查询到有效的原始链接，写入缓存
	key = fmt.Sprintf(consts.GotoShortLinkKey, fullShortUrl)
	err = q.rdb.SetEx(ctx, key, linkPo.OriginalUrl, toolkit.GetLinkCacheExpiration(linkPo.ValidDate)).Err()
	if err != nil {
		return "", err
	}
	return linkPo.OriginalUrl, err
}

func (q LinkQuery) PageLink(ctx context.Context, param query.PageLink) (res *types.PageResp[query.Link], err error) {
	linkPage, err := q.linkDao.PageLink(ctx, param.Gid, consts.StatusEnable, param.OrderTag, param.Current, param.Size)

	if err != nil {
		return
	}

	res = types.ConvertRecords(linkPage, func(linkDTO dao.LinkDTO) query.Link {
		var queryLink query.Link
		err := copier.Copy(&linkDTO, &queryLink)
		if err != nil {
			return query.Link{}
		}
		return queryLink
	})
	return
}

func (q LinkQuery) ListGroupLinkCount(ctx context.Context, gidList []string) (res []query.GroupLinkCount, err error) {
	linkGidCountDTO, err := q.linkDao.ListGroupLinkCount(ctx, gidList)
	if err != nil {
		return
	}

	res = make([]query.GroupLinkCount, 0)
	for _, dto := range linkGidCountDTO {
		res = append(res, query.GroupLinkCount{
			Gid:   dto.Gid,
			Count: dto.Count,
		})
	}
	return
}
