package link

type Goto struct {
	id           int64
	gid          string
	fullShortUrl string
}

func NewLinkGoto(gid, fullShortUrl string) Goto {
	return Goto{
		gid:          gid,
		fullShortUrl: fullShortUrl,
	}
}
