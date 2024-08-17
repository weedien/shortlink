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

type CreateLinkBatchHandler struct {
	repo link.Repository
}

type CreateLinkBatch struct {
	// 原始链接
	OriginUrls []string
	// 描述
	Descriptions []string
	// 分组ID
	Gid string
	// 创建类型 0:接口创建 1:控制台创建
	CreateType int
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidDateType int
	// 有效期
	ValidDate time.Time
}

func (h CreateLinkBatchHandler) Handle(
	ctx context.Context,
	cmd CreateLinkBatch,
) (resp *valobj.ShortLinkCreateBatchVo, err error) {

	var shortLinkRespList []valobj.ShortLinkCreateVo
	for idx, v := range cmd.OriginUrls {
		link, err := types.NewLink(
			v,
			cmd.Gid,
			cmd.CreateType,
			cmd.ValidDateType,
			cmd.ValidDate,
			cmd.Descriptions[idx],
		)
		if err != nil {
			return
		}
		linkGoto := entity.NewLinkGoto(cmd.Gid, link.ShortUrl())
		linkAggregate := aggregate.NewCreateLinkAggregate(link, linkGoto)

		if err = h.repo.CreateLink(ctx, linkAggregate); err != nil {
			return
		}

		r := valobj.ShortLinkCreateVo{
			FullShortUrl: link.ShortUrl(),
			OriginalUrl:  v,
			Gid:          cmd.Gid,
		}
		shortLinkRespList = append(shortLinkRespList, r)
	}

	resp = &valobj.ShortLinkCreateBatchVo{
		SuccessCount:          len(shortLinkRespList),
		ShortLinkCreateVoList: shortLinkRespList,
	}
	return
}
