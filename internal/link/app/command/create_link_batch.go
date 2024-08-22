package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain"
	"shortlink/internal/link/domain/aggregate"
	"shortlink/internal/link/domain/entity"
	"shortlink/internal/link/domain/valobj"
	"time"
)

type createLinkBatchHandler struct {
	repo domain.LinkRepository
}

type CreateLinkBatchHandler decorator.CommandHandler[CreateLinkBatch]

func NewCreateLinkBatchHandler(
	repo domain.LinkRepository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) CreateLinkBatchHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[CreateLinkBatch](
		createLinkBatchHandler{repo: repo},
		logger,
		metricsClient,
	)
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
	// 执行结果
	result *valobj.ShortLinkCreateBatchVo
}

func (c CreateLinkBatch) ExecutionResult() *valobj.ShortLinkCreateBatchVo {
	return c.result
}

func (h createLinkBatchHandler) Handle(
	ctx context.Context,
	cmd CreateLinkBatch,
) error {

	var shortLinkRespList []valobj.ShortLinkCreateVo
	for idx, v := range cmd.OriginUrls {
		linkEntity, err := entity.NewLink(
			v,
			cmd.Gid,
			cmd.CreateType,
			cmd.ValidDateType,
			cmd.ValidDate,
			cmd.Descriptions[idx],
		)
		if err != nil {
			return err
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
			return err
		}

		linkGoto := entity.NewLinkGoto(cmd.Gid, linkEntity.FullShortUrl())
		linkAggregate := aggregate.NewCreateLinkAggregate(linkEntity, linkGoto)

		if err = h.repo.CreateLink(ctx, linkAggregate); err != nil {
			return err
		}

		r := valobj.ShortLinkCreateVo{
			FullShortUrl: linkEntity.FullShortUrl(),
			OriginalUrl:  v,
			Gid:          cmd.Gid,
		}
		shortLinkRespList = append(shortLinkRespList, r)
	}

	cmd.result = &valobj.ShortLinkCreateBatchVo{
		SuccessCount: len(shortLinkRespList),
		LinkInfos:    shortLinkRespList,
	}
	return nil
}
