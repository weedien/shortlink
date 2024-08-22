package assembler

import (
	"shortlink/internal/link/adapter/po"
	"shortlink/internal/link/domain/entity"
)

type LinkRecycleBinConverter struct {
}

func (c LinkRecycleBinConverter) LinkEntityToPo(e entity.Link) po.Link {
	return po.Link{}
}

func (c LinkRecycleBinConverter) LinkPoToEntity(p po.Link) entity.Link {
	return entity.Link{}
}
