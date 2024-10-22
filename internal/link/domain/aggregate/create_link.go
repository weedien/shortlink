package aggregate

import (
	"shortlink/internal/link/domain/link"
)

type CreateLinkAggregate struct {
	ShortLink     link.Link
	ShortLinkGoto link.Goto
}

func NewCreateLinkAggregate(shortLink link.Link, shortLinkGoto link.Goto) CreateLinkAggregate {
	return CreateLinkAggregate{
		ShortLink:     shortLink,
		ShortLinkGoto: shortLinkGoto,
	}
}
