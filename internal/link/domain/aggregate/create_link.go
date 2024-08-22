package aggregate

import (
	"shortlink/internal/link/domain/entity"
)

type CreateLinkAggregate struct {
	ShortLink     entity.Link
	ShortLinkGoto entity.LinkGoto
}

func NewCreateLinkAggregate(shortLink entity.Link, shortLinkGoto entity.LinkGoto) CreateLinkAggregate {
	return CreateLinkAggregate{
		ShortLink:     shortLink,
		ShortLinkGoto: shortLinkGoto,
	}
}
