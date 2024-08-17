package command

import (
	"context"
	"shortlink/common/consts"
	"shortlink/common/types"
	"shortlink/internal/domain/link"
	"time"
)

type UpdateLinkHandler struct {
	repo link.Repository
}

type UpdateLink struct {
	// 完整短链接
	FullShortLink string
	// 原始链接
	OriginalUrl string
	// 原始分组ID
	OriginalGid string
	// 分组ID
	Gid string
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidDateType int
	// 有效期
	ValidDate time.Time
	// 描述
	Description string
}

func (h UpdateLinkHandler) Handle(ctx context.Context, cmd UpdateLink) (err error) {
	return h.repo.UpdateLink(
		ctx,
		types.LinkID{
			FullShortUrl: cmd.FullShortLink,
			Gid:          cmd.OriginalGid,
		},
		consts.StatusEnable,
		func(ctx context.Context, link *types.Link) (*types.Link, error) {
			if err := link.VerifyWhiteList(cmd.OriginalUrl); err != nil {
				return nil, err
			}
			link.SetGid(cmd.Gid).SetOriginalUrl(cmd.OriginalUrl).
				SetDesc(cmd.Description).SetValidDateType(cmd.ValidDateType).
				SetValidDate(cmd.ValidDate)
			return link, nil
		},
	)
}
