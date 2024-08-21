package query

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
)

type listGroupCountHandler struct {
	readModel ListGroupCountReadModel
}

type ListGroupCountHandler decorator.QueryHandler[[]string, []GroupLinkCount]

func NewListGroupCountHandler(
	readModel ListGroupCountReadModel,
	logger *slog.Logger,
	metrics metrics.Client,
) ListGroupCountHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[[]string, []GroupLinkCount](
		listGroupCountHandler{readModel: readModel},
		logger,
		metrics,
	)
}

type ListGroupCountReadModel interface {
	ListGroupLinkCount(ctx context.Context, gidList []string) ([]GroupLinkCount, error)
}

func (h listGroupCountHandler) Handle(ctx context.Context, gidList []string) ([]GroupLinkCount, error) {
	return h.readModel.ListGroupLinkCount(ctx, gidList)
}
