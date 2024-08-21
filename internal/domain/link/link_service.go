package link

import (
	"github.com/google/uuid"
	"shortlink/common/config"
	"shortlink/common/error_no"
	"shortlink/common/toolkit"
)

type LinkService struct {
}

func (s LinkService) VerifyWhiteList(originUrl string) error {
	if !config.EnableWhiteList.Bool() {
		return nil
	}
	domain := toolkit.ExtractDomain(originUrl)
	if domain == "" {
		return error_no.NewServiceError(error_no.OriginalUrlMisspelled)
	}
	whiteList := config.DomainWhiteList.Array()
	for _, v := range whiteList {
		if domain == v {
			return nil
		}
	}
	return error_no.NewServiceErrorWithMsg(
		error_no.InvalidDomain,
		"为避免恶意攻击,只支持生成以下网站跳转链接: "+config.DomainWhiteListNames.String(),
	)
}

func (s LinkService) GenerateSuffix(originalUrl string) string {
	originalUrl += uuid.New().String()
	return toolkit.HashToBase62(originalUrl)
}

func (s LinkService) GenerateFunc(originalUrl string) func() string {
	return func() string {
		return s.GenerateSuffix(originalUrl)
	}
}
