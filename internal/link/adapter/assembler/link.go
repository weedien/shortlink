package assembler

import (
	"shortlink/internal/link/adapter/po"
	"shortlink/internal/link/domain/link"
)

type LinkConverter struct {
}

func (s *LinkConverter) LinkModelToQuery() {

}

func (s *LinkConverter) LinkEntityToPo(entity link.Link) po.Link {
	return po.Link{}
}

func (s *LinkConverter) LinkPoToEntity(po po.Link) *link.Link {
	return &link.Link{}
}

func (s *LinkConverter) LinkGotoEntityToPo(entity link.Link) po.LinkGoto {
	return po.LinkGoto{}
}

func (s *LinkConverter) LinkGotoPoToEntity(po po.LinkGoto) link.Goto {
	return link.Goto{}
}
