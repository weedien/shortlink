package converter

import (
	"shortlink/internal/infra/persistence/po"
)

type LinkRecycleBinConverter struct {
}

func (c LinkRecycleBinConverter) LinkEntityToPo(e types.ShortLink) po.Link {
	return po.Link{}
}

func (c LinkRecycleBinConverter) LinkPoToEntity(p po.Link) types.ShortLink {
	return types.ShortLink{}
}
