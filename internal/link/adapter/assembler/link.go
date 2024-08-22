package assembler

import (
	po2 "shortlink/internal/link/adapter/po"
	"shortlink/internal/link/domain/entity"
)

type LinkConverter struct {
}

func (s *LinkConverter) LinkModelToQuery() {

}

func (s *LinkConverter) LinkEntityToPo(entity entity.Link) po2.Link {
	return po2.Link{}
}

func (s *LinkConverter) LinkPoToEntity(po po2.Link) *entity.Link {
	return &entity.Link{}
}

func (s *LinkConverter) LinkGotoEntityToPo(entity entity.LinkGoto) po2.LinkGoto {
	return po2.LinkGoto{}
}

func (s *LinkConverter) LinkGotoPoToEntity(po po2.LinkGoto) entity.LinkGoto {
	return entity.LinkGoto{}
}
