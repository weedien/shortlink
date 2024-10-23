package link

import (
	"errors"
	"time"
)

const (
	// ValidTypePermanent 永久有效
	ValidTypePermanent = 0
	// ValidTypeTemporary 临时有效
	ValidTypeTemporary = 1
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
	startDate  time.Time
	endDate    time.Time
	hasExpired bool
}

func (v ValidDate) ValidType() int {
	return v.validType
}

func (v ValidDate) HasExpired() bool {
	return v.hasExpired
}

func (v ValidDate) StartDate() time.Time {
	return v.startDate
}

func (v ValidDate) EndDate() time.Time {
	return v.endDate
}

func NewValidDate(validType int, startTime, endTime time.Time) (*ValidDate, error) {
	if validType != ValidTypePermanent && validType != ValidTypeTemporary {
		return &ValidDate{}, errors.New("invalid validType")
	}
	if validType != ValidTypePermanent && endTime.Before(startTime) {
		return &ValidDate{}, errors.New("endDate should be after startDate")
	}
	var hasExpired bool
	if validType == ValidTypePermanent || endTime.After(time.Now()) {
		hasExpired = false
	} else {
		hasExpired = true
	}

	return &ValidDate{
		validType:  validType,
		startDate:  startTime,
		endDate:    endTime,
		hasExpired: hasExpired,
	}, nil
}

func (v ValidDate) isValid() bool {
	return v.validType == ValidTypePermanent || v.endDate.After(time.Now())
}

func (v ValidDate) Expiration() time.Duration {
	if v.validType == ValidTypePermanent {
		return 0
	}
	return v.endDate.Sub(time.Now())
}

func (v ValidDate) NeverExpire() bool {
	return v.validType == ValidTypePermanent
}

func (v ValidDate) StartTime() time.Time {
	return v.startDate
}

func (v ValidDate) EndTime() time.Time {
	return v.endDate
}
