package query

import (
	"context"
	"shortlink/common/consts"
	"shortlink/common/types"
)

type PageDisabledLinkHandler struct {
	readModel PageDisableLinkReadModel
}

type PageRecycleBinParam struct {
	types.PageReq
	GidList      []string `json:"gidList"`
	EnableStatus int      `json:"enableStatus"`
}

type PageDisableLinkReadModel interface {
	PageDisabledLink(ctx context.Context, param PageRecycleBinParam) (*types.PageResp[LinkQueryDTO], error)
}

func (h PageDisabledLinkHandler) Handle(ctx context.Context, gidList []string, page, size int) (*types.PageResp[LinkQueryDTO], error) {
	param := PageRecycleBinParam{
		PageReq: types.PageReq{
			Current: page,
			Size:    size,
		},
		GidList:      gidList,
		EnableStatus: consts.StatusDisable,
	}

	return h.readModel.PageDisabledLink(ctx, param)
}
