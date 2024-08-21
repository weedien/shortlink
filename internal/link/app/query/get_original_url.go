package query

import (
	"context"
	"log/slog"
	"shortlink/internal/common/base_event"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/app/event"
	"shortlink/internal/link/domain/valobj"
)

type getOriginalUrlHandler struct {
	readModel GetOriginalUrlReadModel
	eventBus  base_event.AppEventBus
}

type GetOriginalUrlHandler decorator.QueryHandler[GetOriginalUrl, string]

func NewGetOriginalUrlHandler(
	readModel GetOriginalUrlReadModel,
	eventBus base_event.AppEventBus,
	logger *slog.Logger,
	metrics metrics.Client,
) GetOriginalUrlHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[GetOriginalUrl, string](
		getOriginalUrlHandler{readModel: readModel, eventBus: eventBus},
		logger,
		metrics,
	)
}

type GetOriginalUrl struct {
	FullShortUrl string
	RecordInfo   valobj.ShortLinkStatsRecordVo
}

type GetOriginalUrlReadModel interface {
	// GetOriginalUrlByShortUrl 通过短链接获取原始链接
	GetOriginalUrlByShortUrl(ctx context.Context, fullShortUrl string) (string, error)
}

func (h getOriginalUrlHandler) Handle(ctx context.Context, q GetOriginalUrl) (string, error) {

	// $$ 发布事件 LinkAccessEvent
	e := event.NewLinkAccessedEvent(q.RecordInfo)
	h.eventBus.Publish(ctx, e)

	return h.readModel.GetOriginalUrlByShortUrl(ctx, q.FullShortUrl)
}
