package aggregate

import (
	"shortlink/common/types"
	"time"
)

type UpdateShortLinkAggregate struct {
	OldShortLink types.Link
	NewShortLink types.Link
}

func NewUpdateShortLinkAggregate(old types.ShortLink, new types.ShortLink) UpdateShortLinkAggregate {
	return UpdateShortLinkAggregate{
		OldShortLink: old,
		NewShortLink: new,
	}
}

func (a UpdateShortLinkAggregate) GidChanged() bool {
	return a.OldShortLink.Gid != a.NewShortLink.Gid
}

func (a UpdateShortLinkAggregate) ShouldUpdateCache() bool {
	o := a.OldShortLink
	n := a.NewShortLink

	return o.ValidDateType != n.ValidDateType || o.ValidDate != n.ValidDate || o.OriginalUrl != n.OriginalUrl
}

// ValidStatusChanged 短链接从无效变为有效
func (a UpdateShortLinkAggregate) ValidStatusChanged() bool {
	now := time.Now()
	o := a.OldShortLink
	n := a.NewShortLink
	return o.ValidDate.Before(now) && (n.ValidDate.After(now) || n.ValidDateType == 0)
}
