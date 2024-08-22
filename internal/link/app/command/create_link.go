package command

import (
	"context"
	"fmt"
	"log/slog"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/lock"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain"
	"shortlink/internal/link/domain/aggregate"
	"shortlink/internal/link/domain/entity"
	"shortlink/internal/link/domain/valobj"
	"time"
)

type createLinkHandler struct {
	repo   domain.LinkRepository
	locker lock.DistributedLock
}

type CreateLink struct {
	// 原始链接
	OriginalUrl string
	// 分组ID
	Gid string
	// 创建类型 0:接口创建 1:控制台创建
	CreateType int
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidDateType int
	// 有效期
	ValidDate time.Time
	// 描述
	Description string
	// 是否加锁
	WithLock bool
	// 执行结果
	result *valobj.ShortLinkCreateVo
}

func (c CreateLink) ExecutionResult() *valobj.ShortLinkCreateVo {
	return c.result
}

type CreateLinkHandler decorator.CommandHandler[CreateLink]

func NewCreateLinkHandler(
	repo domain.LinkRepository,
	locker lock.DistributedLock,
	logger *slog.Logger,
	metrics metrics.Client,
) CreateLinkHandler {
	if repo == nil {
		panic("nil repo")
	}
	if locker == nil {
		panic("nil locker")
	}

	return decorator.ApplyCommandDecorators[CreateLink](
		createLinkHandler{repo: repo, locker: locker},
		logger,
		metrics,
	)
}

func (h createLinkHandler) Handle(
	ctx context.Context,
	cmd CreateLink,
) (err error) {

	// 获取分布式锁
	if cmd.WithLock {
		lockKey := fmt.Sprintf(constant.ShortLinkCreateLockKey, cmd.OriginalUrl)
		if _, err = h.locker.Acquire(ctx, lockKey, time.Second); err != nil {
			return err
		}
		defer func() {
			if err := h.locker.Release(ctx, lockKey); err != nil {
				err = fmt.Errorf("释放锁异常: %w", err)
			}
		}()
	}

	// 创建短链接实体
	linkEntity, err := entity.NewLink(
		cmd.OriginalUrl,
		cmd.Gid,
		cmd.CreateType,
		cmd.ValidDateType,
		cmd.ValidDate,
		cmd.Description,
	)
	if err != nil {
		return
	}

	// 生成唯一短链接
	err = linkEntity.GenUniqueShortUri(10, func(shortUri string) bool {
		exists, err := h.repo.ShortUriExists(ctx, shortUri)
		if err != nil {
			return true
		}
		return exists
	})
	if err != nil {
		return
	}

	// 持久化短链接
	linkGoto := entity.NewLinkGoto(cmd.Gid, linkEntity.FullShortUrl())
	linkAggregate := aggregate.NewCreateLinkAggregate(linkEntity, linkGoto)

	if err = h.repo.CreateLink(ctx, linkAggregate); err != nil {
		return
	}

	cmd.result = &valobj.ShortLinkCreateVo{
		FullShortUrl: linkEntity.FullShortUrl(),
		OriginalUrl:  cmd.OriginalUrl,
		Gid:          cmd.Gid,
	}
	return
}
