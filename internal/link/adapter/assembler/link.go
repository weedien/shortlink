package assembler

import (
	"shortlink/internal/link/adapter/po"
	"shortlink/internal/link/domain/link"
)

type LinkAssembler struct {
	linkFactory *link.Factory
}

func (a *LinkAssembler) LinkEntityToLinkPo(lk *link.Link) *po.Link {
	return &po.Link{
		Gid:         lk.Gid(),
		ShortUri:    lk.ShortUri(),
		OriginalUrl: lk.OriginalUrl(),
		Favicon:     lk.Favicon(),
		Status:      lk.Status(),
		CreateType:  lk.CreateType(),
		ValidType:   lk.ValidDate().ValidType(),
		StartDate:   lk.ValidDate().StartDate(),
		EndDate:     lk.ValidDate().EndDate(),
		Desc:        lk.Desc(),
	}
}

func (a *LinkAssembler) LinkPoToLinkEntity(po *po.Link) *link.Link {
	validDate, err := link.NewValidDate(po.ValidType, po.StartDate, po.EndDate)
	if err != nil {
		return nil
	}

	lk := &link.Link{}
	if lk, err = a.linkFactory.NewLinkFromDB(
		po.ID,
		po.Gid,
		po.ShortUri,
		po.OriginalUrl,
		po.Status,
		po.CreateType,
		po.Favicon,
		po.Desc,
		validDate,
	); err != nil {
		return nil
	}
	return lk
}

func (a *LinkAssembler) LinkEntityToLinkGotoPo(lk *link.Link) *po.LinkGoto {
	return &po.LinkGoto{
		Gid:      lk.Gid(),
		ShortUri: lk.ShortUri(),
	}
}
