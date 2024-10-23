package query

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"shortlink/internal/common/base_event"
	"shortlink/internal/common/cache"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/error_no"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain/event"
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
	distributedCache cache.DistributedCache,
	logger *slog.Logger,
	metrics metrics.Client,
) GetOriginalUrlHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[GetOriginalUrl, string](
		getOriginalUrlHandler{readModel: readModel, eventBus: eventBus, distributedCache: distributedCache},
		logger,
		metrics,
	)
}

type GetOriginalUrl struct {
	ShortUri      string
	UserVisitInfo event.UserVisitInfo
}

type GetOriginalUrlReadModel interface {
	GetLink(ctx context.Context, shortUri string) (*link.Link, error)
}

func (h getOriginalUrlHandler) Handle(ctx context.Context, q GetOriginalUrl) (res string, err error) {

	// $$ 发布事件 UserVisitEvent
	e := event.NewUserVisitEvent(q.UserVisitInfo)
	h.eventBus.Publish(ctx, e)

	fetchFn := func() (res interface{}, err error) {
		lk := &link.Link{}
		if lk, err = h.readModel.GetLink(ctx, q.ShortUri); err != nil || lk == nil {
			return
		}
		originalUrl := lk.OriginalUrl()
		if originalUrl == "" {
			err = error_no.LinkNotExists
			return
		}
		switch lk.Status() {
		case link.StatusActive:
			res = link.NewCacheValue(lk)
			return
		case link.StatusExpired:
			err = error_no.LinkExpired
			return
		case link.StatusForbidden:
			err = error_no.LinkForbidden
			return
		case link.StatusReserved:
			err = error_no.LinkReserved
			return
		default:
			err = error_no.LinkNotExists
			return
		}
	}

	var result interface{}
	result, err = h.distributedCache.SafeGetWithCacheCheckFilter(
		ctx,
		constant.GotoLinkKey+q.ShortUri,
		reflect.TypeOf(link.CacheValue{}),
		fetchFn,
		constant.NeverExpire,
		cache.ShortUriCreateBloomFilter,
		q.ShortUri,
		constant.GotoIsNullLinkKey+q.ShortUri,
	)
	if err != nil {
		if errors.Is(err, error_no.RedisKeyNotExist) {
			err = error_no.LinkNotExists
		}
		return
	}

	if cacheValue, ok := result.(*link.CacheValue); ok {
		res = cacheValue.OriginalUrl
	}

	return
}
