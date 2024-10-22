package link

import (
	"errors"
	"time"
)

const (
	// ValidTypePermanent 永久有效
	ValidTypePermanent = -1
	// ValidTypeTemporary 临时有效
	ValidTypeTemporary = 0
)

func (lk Link) validateDate() error {
	if lk.validDate.isValid() {
		return nil
	}
	return errors.New("link has expired")
}

// ValidDate 有效期
type ValidDate struct {
	validType  int
	startTime  time.Time
	endTime    time.Time
	hasExpired bool
}

func NewValidDate(validType int, startTime, endTime time.Time) (*ValidDate, error) {
	if validType != ValidTypePermanent && validType != ValidTypeTemporary {
		return &ValidDate{}, errors.New("invalid validType")
	}
	if endTime.Before(startTime) {
		return &ValidDate{}, errors.New("endTime should be after startTime")
	}
	var hasExpired bool
	if validType == ValidTypePermanent || endTime.After(time.Now()) {
		hasExpired = false
	} else {
		hasExpired = true
	}

	return &ValidDate{
		validType:  validType,
		startTime:  startTime,
		endTime:    endTime,
		hasExpired: hasExpired,
	}, nil
}

func (v ValidDate) isValid() bool {
	return v.validType == ValidTypePermanent || v.endTime.After(time.Now())
}

func (v ValidDate) Expiration() time.Duration {
	return v.endTime.Sub(time.Now())
}

func (v ValidDate) NeverExpire() bool {
	return v.validType == ValidTypePermanent
}

func (v ValidDate) StartTime() time.Time {
	return v.startTime
}

func (v ValidDate) EndTime() time.Time {
	return v.endTime
}
