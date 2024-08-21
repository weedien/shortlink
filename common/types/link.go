package types

import (
	"errors"
	"github.com/google/uuid"
	"shortlink/common/config"
	"shortlink/common/consts"
	"shortlink/common/toolkit"
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
	// 白名单校验
	err = verifyWhiteList(originalUrl)
	if err != nil {
		return
	}

	// 有效期需要大于当前时间
	if validDate.Before(time.Now()) {
		err = errors.New("有效期不能小于当前时间")
		return
	}

	domain := config.ShortLinkDomain.String()
	favicon := toolkit.GetFaviconWithDefault(originalUrl, config.DefaultFavicon.String())

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

	return
}

func verifyWhiteList(originUrl string) error {
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

func (l Link) FullShortUrl() string {
	return l.fullShortUrl
}

// GenUniqueShortUri 生成唯一短链接
func (l Link) GenUniqueShortUri(attempts int, ifExistsFunc func(string) bool) error {
	for i := 0; i < attempts; i++ {
		shortUri := toolkit.HashToBase62(l.originalUrl)
		if !ifExistsFunc(shortUri) {
			l.shortUri = "https://" + l.domain + "/" + shortUri
			return nil
		}
	}
	return errors.New("多次尝试生成唯一短链接失败")
}

func (l Link) genShortUrl() {
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
