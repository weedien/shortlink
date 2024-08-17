package command

import (
	"context"
	"shortlink/common/types"
	"shortlink/internal/domain/link"
	"shortlink/internal/domain/link/aggregate"
	"shortlink/internal/domain/link/entity"
	"shortlink/internal/domain/link/valobj"
	"time"
)

type CreateLinkHandler struct {
	repo link.Repository
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
}

func (h CreateLinkHandler) Handle(
	ctx context.Context,
	cmd CreateLink,
) (resp *valobj.ShortLinkCreateVo, err error) {

	link, err := types.NewLink(
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
	linkGoto := entity.NewLinkGoto(cmd.Gid, link.ShortUrl())
	linkAggregate := aggregate.NewCreateLinkAggregate(link, linkGoto)

	if cmd.WithLock {
		if err = h.repo.CreateLinkWithLock(ctx, linkAggregate); err != nil {
			return
		}
	} else {
		if err = h.repo.CreateLink(ctx, linkAggregate); err != nil {
			return
		}
	}

	resp = &valobj.ShortLinkCreateVo{
		FullShortUrl: link.ShortUrl(),
		OriginalUrl:  cmd.OriginalUrl,
		Gid:          cmd.Gid,
	}
	return
}
