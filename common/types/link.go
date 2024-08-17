package types

import (
	"errors"
	"github.com/google/uuid"
	"shortlink/common/consts"
	"shortlink/common/toolkit"
	"shortlink/config"
	"time"
)

// LinkID 短链接唯一标识
type LinkID struct {
	// 分组ID
	Gid string
	// 完整短链接
	FullShortUrl string
}

type Link struct {
	id            int
	domain        string
	shortUri      string
	fullShortUrl  string
	originalUrl   string
	clickNum      int
	gid           string
	enableStatus  int
	createType    int
	validDateType int
	validDate     time.Time
	desc          string
	favicon       string
	totalPv       int
	totalUv       int
	totalUip      int
	todayPv       int
	todayUv       int
	todayUip      int
}

func NewLink(
	originalUrl string,
	gid string,
	createType int,
	validDateType int,
	validDate time.Time,
	desc string,
) (link Link, err error) {

	domain := config.ShortLinkDomain.String()
	favicon := toolkit.GetFaviconWithDefault(originalUrl, config.DefaultFavicon.String())

	// TODO 校验validDate是否有效

	link = Link{
		originalUrl:   originalUrl,
		domain:        domain,
		gid:           gid,
		createType:    createType,
		validDateType: validDateType,
		validDate:     validDate,
		desc:          desc,
		favicon:       favicon,
		enableStatus:  consts.StatusEnable,
	}

	// 白名单校验
	err = link.VerifyWhiteList(originalUrl)

	// 生成短链接
	link.generateShortUrl()

	return
}

func (l Link) VerifyWhiteList(originUrl string) error {
	if !config.EnableWhiteList.Bool() {
		return nil
	}
	domain := toolkit.ExtractDomain(originUrl)
	if domain == "" {
		return errors.New("原始链接拼写错误")
	}
	whiteList := config.DomainWhiteList.Array()
	for _, v := range whiteList {
		if domain == v {
			return nil
		}
	}
	return errors.New("为避免恶意攻击,只支持生成以下网站跳转链接: " + config.DomainWhiteListNames.String())
}

func (l Link) Enable() {
	l.enableStatus = consts.StatusEnable
}
func (l Link) Disable() {
	l.enableStatus = consts.StatusDisable
}

func (l Link) ShortUrl() string {
	return l.fullShortUrl
}

func (l Link) generateShortUrl() {
	originalUrl := l.originalUrl
	originalUrl += uuid.New().String()
	l.fullShortUrl = "https://" + l.domain + "/" + toolkit.HashToBase62(originalUrl)
}

//func (l Link) GenerateFunc(originalUrl string) func() string {
//	return func() string {
//		return l.generateSuffix(originalUrl)
//	}
//}

func (l Link) SetGid(gid string) Link {
	l.gid = gid
	return l
}

func (l Link) SetDesc(desc string) Link {
	l.desc = desc
	return l
}

func (l Link) SetValidDate(validDate time.Time) Link {
	l.validDate = validDate
	return l
}

func (l Link) SetValidDateType(validDateType int) Link {
	l.validDateType = validDateType
	return l
}

func (l Link) SetCreateType(createType int) Link {
	l.createType = createType
	return l
}

func (l Link) SetOriginalUrl(originalUrl string) Link {
	l.originalUrl = originalUrl
	return l
}

func (l Link) SetFavicon(favicon string) Link {
	l.favicon = favicon
	return l
}

func (l Link) ValidDate() time.Time {
	return l.validDate
}
