package converter

import (
	"shortlink/common/types"
	"shortlink/internal/domain/link/entity"
	po2 "shortlink/internal/infra/persistence/po"
)

type LinkConverter struct {
}

func (s *LinkConverter) LinkModelToQuery() {

}

func (s *LinkConverter) LinkEntityToPo(entity *types.Link) po2.Link {
	return po2.Link{}
}

func (s *LinkConverter) LinkPoToEntity(po po2.Link) *types.Link {
	return &types.Link{}
}

func (s *LinkConverter) LinkGotoEntityToPo(entity entity.LinkGoto) po2.LinkGoto {
	return po2.LinkGoto{}
}

func (s *LinkConverter) LinkGotoPoToEntity(po po2.LinkGoto) entity.LinkGoto {
	return entity.LinkGoto{}
}
