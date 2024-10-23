package link

import (
	"errors"
	"fmt"
	"shortlink/internal/common/toolkit"
	"strings"
	"time"
)

type FactoryConfig struct {
	Domain         string
	UseSSL         bool
	Whitelist      []string
	DefaultFavicon string
	MaxAttempts    int // 生成唯一短链接的最大尝试次数
}

func (f FactoryConfig) Validate() error {
	var err error

	if f.Domain == "" {
		err = errors.Join(err, errors.New("domain should not be empty"))
	}
	if !toolkit.IsValidDomain(f.Domain) {
		err = errors.Join(err, errors.New("domain should be valid, but is "+f.Domain))
	}
	for _, v := range f.Whitelist {
		if !toolkit.IsValidDomain(v) {
			err = errors.Join(err, errors.New("whitelist domain should be valid, but is "+v))
		}
	}
	if f.DefaultFavicon == "" {
		if !toolkit.IsValidUrl(f.DefaultFavicon) {
			err = errors.Join(err, errors.New("default favicon should be valid url"))
		}
	}
	if f.MaxAttempts < 1 {
		err = errors.Join(
			err,
			errors.New("MaxAttempts should be greater than 1, but is "+fmt.Sprint(f.MaxAttempts)),
		)
	}

	return err
}

// Factory 全局唯一短链接工厂
type Factory struct {
	fc FactoryConfig
}

func NewFactory(fc FactoryConfig) (*Factory, error) {
	if err := fc.Validate(); err != nil {
		return &Factory{}, errors.Join(err, errors.New("invalid config passed to link factory"))
	}

	return &Factory{fc: fc}, nil
}

func (f Factory) NewAvailableLink(
	originalUrl string,
	gid string,
	createType int,
	validType int,
	validEndDate time.Time,
	desc string,
	ifExistsFunc func(string) (bool, error),
) (lk *Link, err error) {

	// 白名单校验
	if err = f.verifyWhiteList(originalUrl); err != nil {
		return nil, err
	}

	// 短链接
	var shortUri string
	if shortUri, err = f.genUniqueShortUri(originalUrl, f.fc.MaxAttempts, ifExistsFunc); err != nil {
		return nil, err
	}

	// 完整短链接
	var fullShortUrl string
	if f.fc.UseSSL {
		fullShortUrl = fmt.Sprintf("https://%s/%s", f.fc.Domain, shortUri)
	} else {
		fullShortUrl = fmt.Sprintf("http://%s/%s", f.fc.Domain, shortUri)
	}

	// 原网址图标
	var favicon string
	if favicon, _ = toolkit.GetFavicon(originalUrl); favicon == "" {
		favicon = f.fc.DefaultFavicon
	}

	// 有效期
	var validDate *ValidDate
	if validDate, err = NewValidDate(validType, time.Now(), validEndDate); err != nil {
		return nil, err
	}

	return &Link{
		domain:       f.fc.Domain,
		shortUri:     shortUri,
		fullShortUrl: fullShortUrl,
		originalUrl:  originalUrl,
		gid:          gid,
		status:       StatusActive,
		createType:   createType,
		validDate:    validDate,
		desc:         desc,
		favicon:      favicon,
	}, nil
}

func (f Factory) NewLinkFromDB(
	id int,
	gid string,
	shortUri string,
	originalUrl string,
	status string,
	createType int,
	favicon string,
	desc string,
	validDate *ValidDate,
) (*Link, error) {
	// 完整短链接
	var fullShortUrl string
	if f.fc.UseSSL {
		fullShortUrl = fmt.Sprintf("https://%s/%s", f.fc.Domain, shortUri)
	} else {
		fullShortUrl = fmt.Sprintf("http://%s/%s", f.fc.Domain, shortUri)
	}

	return &Link{
		id:           id,
		domain:       f.fc.Domain,
		shortUri:     shortUri,
		fullShortUrl: fullShortUrl,
		originalUrl:  originalUrl,
		gid:          gid,
		status:       status,
		createType:   createType,
		desc:         desc,
		favicon:      favicon,
		validDate:    validDate,
	}, nil

}

//func (f Factory) NewAvailableLinkWithLock(
//	originalUrl string,
//	gid string,
//	createType int,
//	validType int,
//	validEndDate time.Time,
//	desc string,
//	ifExistsFunc func(string) (bool, error),
//	lockFunc func() error,
//) (lk *Link, err error) {
//
//	if err = lockFunc(); err != nil {
//		return nil, err
//	}
//	return f.NewAvailableLink(originalUrl, gid, createType, validType, validEndDate, desc, ifExistsFunc)
//}

func (f Factory) genUniqueShortUri(
	originalUrl string,
	maxAttempts int,
	ifExistsFunc func(string) (bool, error),
) (shortUri string, err error) {
	for i := 0; i < maxAttempts; i++ {
		shortUri = toolkit.HashToBase62(originalUrl)
		var exists bool
		if exists, err = ifExistsFunc(shortUri); err != nil {
			return "", err
		}
		if !exists {
			return shortUri, nil
		}
	}
	return "", errors.New("多次尝试生成唯一短链接失败")
}

func (f Factory) verifyWhiteList(originUrl string) error {
	whitelist := f.fc.Whitelist
	if whitelist == nil || len(whitelist) == 0 {
		return nil
	}

	domain := toolkit.ExtractDomain(originUrl)
	if domain == "" {
		return errors.New("invalid originUrl: " + originUrl)
	}
	for _, v := range whitelist {
		if domain == v {
			return nil
		}
	}
	return errors.New("白名单只支持跳转以下链接: " + strings.Join(whitelist, ","))
}
