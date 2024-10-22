package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain"
	"shortlink/internal/link/domain/link"
	"time"
)

type createLinkBatchHandler struct {
	repo        domain.LinkRepository
	linkFactory *link.Factory
}

type CreateLinkBatchHandler decorator.CommandHandler[CreateLinkBatch]

func NewCreateLinkBatchHandler(
	repo domain.LinkRepository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) CreateLinkBatchHandler {

	return decorator.ApplyCommandDecorators[CreateLinkBatch](
		createLinkBatchHandler{repo: repo},
		logger,
		metricsClient,
	)
}

type CreateLinkBatch struct {
	// 原始链接
	OriginalUrls []string
	// 描述
	Descs []string
	// 分组ID
	Gid string
	// 创建类型 0:接口创建 1:控制台创建
	CreateType int
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidType int
	// 有效期 - 开始时间
	ValidStartDate time.Time
	// 有效期 - 结束时间
	ValidEndDate time.Time
	// 执行结果
	result *CreateLinkBatchResult
}

type CreateLinkBatchResult struct {
	SuccessCount int
	LinkInfos    []CreateLinkResult
}

func (c CreateLinkBatch) ExecutionResult() *CreateLinkBatchResult {
	return c.result
}

func (h createLinkBatchHandler) Handle(
	ctx context.Context,
	cmd CreateLinkBatch,
) (err error) {

	lks := make([]*link.Link, 0, len(cmd.OriginalUrls))
	linkInfos := make([]CreateLinkResult, 0, len(cmd.OriginalUrls))
	for idx, originalUrl := range cmd.OriginalUrls {
		lk := &link.Link{}
		lk, err = h.linkFactory.NewAvailableLink(
			originalUrl, cmd.Gid, cmd.CreateType, cmd.ValidType, cmd.ValidEndDate, cmd.Descs[idx],
			func(shortUri string) (exists bool, err error) {
				if exists, err = h.repo.ShortUriExists(ctx, shortUri); err != nil {
					return exists, err
				}
				return exists, nil
			},
		)
		if err != nil {
			return err
		}

		lks = append(lks, lk)

		linkInfos[idx] = CreateLinkResult{
			Gid:          cmd.Gid,
			FullShortUrl: lk.FullShortUrl(),
			OriginalUrl:  originalUrl,
		}
	}

	if err = h.repo.CreateLinkBatch(ctx, lks); err != nil {
		return err
	}

	cmd.result = &CreateLinkBatchResult{
		SuccessCount: len(linkInfos),
		LinkInfos:    linkInfos,
	}

	return nil
}
