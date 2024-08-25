package query

import (
	"context"
	"shortlink/internal/user/domain/group"
)

type ListGroupHandler struct {
	readModel        ListGroupReadModel
	shortLinkService ShortLinkService
}

type GroupDto struct {
	Gid            string `json:"gid"`
	Name           string `json:"name"`
	SortOrder      int    `json:"sort_order"`
	ShortLinkCount int    `json:"short_link_count"`
}

type ListGroupReadModel interface {
	ListGroup(ctx context.Context) ([]group.Group, error)
}

func NewListGroupHandler(readModel ListGroupReadModel) ListGroupHandler {
	return ListGroupHandler{readModel: readModel}
}

func (h ListGroupHandler) Handle(ctx context.Context) (res []GroupDto, err error) {
	groups, err := h.readModel.ListGroup(ctx)
	if err != nil {
		return nil, err
	}

	gids := make([]string, len(groups))
	for i, g := range groups {
		gids[i] = g.Gid()
	}

	linkCount, err := h.shortLinkService.ListGroupShortLinkCount(ctx, gids)
	if err != nil {
		return nil, err
	}

	linkCountMap := make(map[string]int, len(linkCount))
	for _, lc := range linkCount {
		linkCountMap[lc.Gid] = lc.ShortLinkCount
	}

	res = make([]GroupDto, len(groups))
	for i, g := range groups {
		res[i] = GroupDto{
			Gid:            g.Gid(),
			Name:           g.Name(),
			SortOrder:      g.SortOrder(),
			ShortLinkCount: linkCountMap[g.Gid()],
		}
	}

	return res, nil
}
