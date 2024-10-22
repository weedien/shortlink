package assembler

import (
	"shortlink/internal/link/adapter/po"
	"shortlink/internal/link/domain/link"
)

type LinkRecycleBinConverter struct {
}

func (c LinkRecycleBinConverter) LinkEntityToPo(e link.Link) po.Link {
	return po.Link{}
}

func (c LinkRecycleBinConverter) LinkPoToEntity(p po.Link) link.Link {
	return link.Link{}
}
