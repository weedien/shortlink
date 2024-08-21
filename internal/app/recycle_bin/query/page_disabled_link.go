package query

import (
	"context"
	"log/slog"
	"shortlink/common/decorator"
	"shortlink/common/types"
)

type pageRecycleBinHandler struct {
	readModel PageRecycleBinReadModel
}

type PageRecycleBinHandler decorator.QueryHandler[PageRecycleBin, *types.PageResp[LinkQueryDTO]]

func NewPageRecycleBinHandler(
	readModel PageRecycleBinReadModel,
	logger *slog.Logger,
	metricsClient metrics.Client,
) PageRecycleBinHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[PageRecycleBin, *types.PageResp[LinkQueryDTO]](
		pageRecycleBinHandler{readModel},
		logger,
		metricsClient,
	)
}

type PageRecycleBin struct {
	types.PageReq
	GidList      []string `json:"gidList"`
	EnableStatus int      `json:"enableStatus"`
}

type PageRecycleBinReadModel interface {
	PageDisabledLink(ctx context.Context, param PageRecycleBin) (*types.PageResp[LinkQueryDTO], error)
}

func (h pageRecycleBinHandler) Handle(ctx context.Context, cmd PageRecycleBin) (*types.PageResp[LinkQueryDTO], error) {
	return h.readModel.PageDisabledLink(ctx, cmd)
}
