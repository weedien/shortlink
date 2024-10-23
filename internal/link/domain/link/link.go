package link

import (
	"errors"
	"shortlink/internal/common/error_no"
	"strings"
	"time"
)

const (
	StatusActive    = "active"    // 激活状态
	StatusExpired   = "expired"   // 已过期
	StatusDisabled  = "disabled"  // 已停用
	StatusForbidden = "forbidden" // 已禁用
	StatusReserved  = "reserved"  // 预留状态
	StatusDeleted   = "deleted"   // 已删除
)

func IsValidStatus(status string) bool {
	return status == StatusActive || status == StatusExpired ||
		status == StatusDisabled || status == StatusForbidden ||
		status == StatusDeleted || status == StatusReserved
}

const (
	CreateByApi = iota
	CreateByConsole
)

// Identifier 短链接唯一标识
//
// 实际上 shortUri 可以确定唯一短链接，而 gid 在分库分表场景下作为分片键，
// 所以 gid 的意义在于避免扩散查询
type Identifier struct {
	// 分组ID
	Gid string
	// 完整短链接
	ShortUri string
}

// Link 短链接
// 由 shortUri 确定唯一短链接
type Link struct {
	id           int
	domain       string
	shortUri     string
	fullShortUrl string
	originalUrl  string
	gid          string
	status       string
	createType   int
	desc         string
	favicon      string
	validDate    *ValidDate
}

func (lk Link) CreateType() int {
	return lk.createType
}

func (lk Link) Gid() string {
	return lk.gid
}

func (lk Link) Favicon() string {
	return lk.favicon
}

func (lk Link) OriginalUrl() string {
	return lk.originalUrl
}

func (lk Link) ShortUri() string {
	return lk.shortUri
}

func (lk Link) Status() string {
	return lk.status
}

func (lk Link) RecoverFromRecycleBin() {
	lk.status = StatusActive
}
func (lk Link) SaveToRecycleBin() {
	lk.status = StatusDisabled
}

func (lk Link) FullShortUrl() string {
	return lk.fullShortUrl
}

func (lk Link) ValidDate() *ValidDate {
	return lk.validDate
}

func (lk Link) Desc() string {
	return lk.desc
}

// Update 更新短链接信息
func (lk Link) Update(
	gid *string,
	originalUrl *string,
	status *string,
	validType *int,
	validEndDate *time.Time,
	desc *string,
) error {
	if gid != nil {
		if strings.TrimSpace(*gid) == "" {
			return errors.New("gid不能为空")
		}
		lk.gid = *gid
	}
	if originalUrl != nil {
		if strings.TrimSpace(*originalUrl) == "" {
			return errors.New("原始链接不能为空")
		}
		lk.originalUrl = *originalUrl
	}
	if status != nil {
		if *status != StatusActive && *status != StatusDisabled {
			return errors.New("invalid status")
		}
		lk.status = *status
	}
	if validType != nil {
		if *validType != ValidTypePermanent && *validType != ValidTypeTemporary {
			return errors.New("invalid validType")
		}
		lk.validDate.validType = *validType
	}
	if validEndDate != nil {
		if validEndDate.Before(time.Now()) {
			return errors.New("endDate should be after startDate")
		}
		lk.validDate.endDate = *validEndDate
	}
	if desc != nil {
		lk.desc = *desc
	}
	return nil
}

type CacheValue struct {
	OriginalUrl    string    `json:"originalUrl"`
	NeverExpire    bool      `json:"neverExpire"`
	ValidStartTime time.Time `json:"validStartTime"`
	ValidEndTime   time.Time `json:"validEndTime"`
	Status         string    `json:"status"`
}

func NewCacheValue(lk *Link) *CacheValue {
	return &CacheValue{
		OriginalUrl:    lk.OriginalUrl(),
		NeverExpire:    lk.ValidDate().NeverExpire(),
		ValidStartTime: lk.ValidDate().StartTime(),
		ValidEndTime:   lk.ValidDate().EndTime(),
		Status:         lk.Status(),
	}
}

func (c CacheValue) Validate() (bool, error) {
	if c.Status == StatusActive {
		if c.NeverExpire {
			return true, nil
		}
		if c.ValidStartTime.Before(time.Now()) && c.ValidEndTime.After(time.Now()) {
			return true, nil
		}
		return false, error_no.LinkExpired
	}
	switch {
	case c.Status == StatusReserved:
		return false, error_no.LinkReserved
	case c.Status == StatusForbidden:
		return false, error_no.LinkForbidden
	case c.Status == StatusDisabled:
		return false, error_no.LinkDisabled
	case c.Status == StatusExpired:
		return false, error_no.LinkExpired
	}
	return false, nil
}

func (c CacheValue) Expiration() time.Duration {
	if c.NeverExpire {
		return 0
	}
	return c.ValidEndTime.Sub(time.Now())
}
