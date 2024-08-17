package query

import "context"

type ListGroupCountHandler struct {
	readModel ListGroupCountReadModel
}

type ListGroupCountReadModel interface {
	ListGroupShortLinkCount(ctx context.Context, gidList []string) ([]GroupLinkCount, error)
}

func (h ListGroupCountHandler) Handle(ctx context.Context, gidList []string) ([]GroupLinkCount, error) {
	return h.readModel.ListGroupShortLinkCount(ctx, gidList)
}
