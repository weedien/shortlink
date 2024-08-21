package aggregate

import (
	"shortlink/internal/common/types"
	"shortlink/internal/link/domain/entity"
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
