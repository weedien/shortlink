package assembler

import (
	"shortlink/internal/common/persistence/po"
	"shortlink/internal/common/types"
)

type LinkRecycleBinConverter struct {
}

func (c LinkRecycleBinConverter) LinkEntityToPo(e types.Link) po.Link {
	return po.Link{}
}

func (c LinkRecycleBinConverter) LinkPoToEntity(p po.Link) types.Link {
	return types.Link{}
}
