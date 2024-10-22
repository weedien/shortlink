package query

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"shortlink/internal/common/base_event"
	"shortlink/internal/common/cache"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/error_no"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/app/event"
	"shortlink/internal/link/domain/link"
)

type getOriginalUrlHandler struct {
	readModel        GetOriginalUrlReadModel
	eventBus         base_event.EventBus
	distributedCache cache.DistributedCache
}

type GetOriginalUrlHandler decorator.QueryHandler[GetOriginalUrl, string]

func NewGetOriginalUrlHandler(
	readModel GetOriginalUrlReadModel,
	eventBus base_event.EventBus,
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
	ShortUri      string
	UserVisitInfo link.UserVisitInfo
}

type GetOriginalUrlReadModel interface {
	GetLinkWithoutStats(ctx context.Context, shortUri string) (*link.Link, error)
}

func (h getOriginalUrlHandler) Handle(ctx context.Context, q GetOriginalUrl) (res string, err error) {

	// $$ 发布事件 UserVisitEvent
	e := event.NewUserVisitEvent(q.UserVisitInfo)
	h.eventBus.Publish(ctx, e)

	fetchFn := func() (res interface{}, err error) {
		lk := &link.Link{}
		if lk, err = h.readModel.GetLinkWithoutStats(ctx, q.ShortUri); err != nil || lk == nil {
			return
		}
		originalUrl := lk.OriginalUrl()
		if originalUrl == "" {
			err = error_no.ShortLinkNotExists
			return
		}
		switch lk.Status() {
		case link.StatusActive:
			res = originalUrl
			return
		case link.StatusExpired:
			err = error_no.ShortLinkExpired
			return
		case link.StatusForbidden:
			err = error_no.ShortLinkForbidden
			return
		case link.StatusReserved:
			err = error_no.ShortLinkReserved
			return
		default:
			err = error_no.ShortLinkNotExists
			return
		}
	}

	var result interface{}
	result, err = h.distributedCache.SafeGetWithCacheCheckFilter(
		ctx,
		constant.GotoShortLinkKey+q.ShortUri,
		reflect.TypeOf(""),
		fetchFn,
		constant.NeverExpire,
		cache.ShortUriCreateBloomFilter,
		fmt.Sprintf(constant.GotoIsNullShortLinkKey, q.ShortUri),
	)
	if err != nil {
		return
	}
	res = result.(string)
	return
}
