package entity

type LinkGoto struct {
	Id           int64
	Gid          string
	FullShortUrl string
}

func NewLinkGoto(gid, fullShortUrl string) LinkGoto {
	return LinkGoto{
		Gid:          gid,
		FullShortUrl: fullShortUrl,
	}
}
