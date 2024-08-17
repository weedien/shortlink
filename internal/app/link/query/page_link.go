package query

import (
	"context"
	"shortlink/common/types"
)

type PageLinkHandler struct {
	readModel PageLinkReadModel
}

type PageLink struct {
	// 分页请求
	types.PageReq
	// 分组ID
	Gid string
	// 排序标识
	OrderTag string
}

type PageLinkReadModel interface {
	PageLink(ctx context.Context, param PageLink) (*types.PageResp[Link], error)
}

func (h PageLinkHandler) Handler(ctx context.Context, param PageLink) (*types.PageResp[Link], error) {
	return h.readModel.PageLink(ctx, param)
}
