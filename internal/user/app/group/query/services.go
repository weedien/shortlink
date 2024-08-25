package query

import "context"

type GroupShortLinkCountDto struct {
	Gid            string `json:"gid"`
	ShortLinkCount int    `json:"short_link_count"`
}

type ShortLinkService interface {
	ListGroupShortLinkCount(ctx context.Context, gids []string) ([]GroupShortLinkCountDto, error)
}
