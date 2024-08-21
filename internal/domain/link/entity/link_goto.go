package entity

type LinkGoto struct {
	id           int64
	gid          string
	fullShortUrl string
}

func NewLinkGoto(gid, fullShortUrl string) LinkGoto {
	return LinkGoto{
		gid:          gid,
		fullShortUrl: fullShortUrl,
	}
}
