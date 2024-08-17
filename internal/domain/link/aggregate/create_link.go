package aggregate

import (
	"shortlink/common/types"
	"shortlink/internal/domain/link/entity"
)

type CreateLinkAggregate struct {
	ShortLink     types.Link
	ShortLinkGoto entity.LinkGoto
}

func NewCreateLinkAggregate(shortLink types.Link, shortLinkGoto entity.LinkGoto) CreateLinkAggregate {
	return CreateLinkAggregate{
		ShortLink:     shortLink,
		ShortLinkGoto: shortLinkGoto,
	}
}
