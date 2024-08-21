package assembler

import (
	"shortlink/internal/common/persistence/po"
	"shortlink/internal/common/types"
	"shortlink/internal/link/domain/entity"
)

type LinkConverter struct {
}

func (s *LinkConverter) LinkModelToQuery() {

}

func (s *LinkConverter) LinkEntityToPo(entity types.Link) po.Link {
	return po.Link{}
}

func (s *LinkConverter) LinkPoToEntity(po po.Link) *types.Link {
	return &types.Link{}
}

func (s *LinkConverter) LinkGotoEntityToPo(entity entity.LinkGoto) po.LinkGoto {
	return po.LinkGoto{}
}

func (s *LinkConverter) LinkGotoPoToEntity(po po.LinkGoto) entity.LinkGoto {
	return entity.LinkGoto{}
}
