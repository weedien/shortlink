package event

import (
	"shortlink/common/base_event"
	"shortlink/internal/domain/link/valobj"
)

type RecordLinkVisitEvent struct {
	base_event.CommonEvent
	RecordInfo valobj.ShortLinkStatsRecordVo
}

func (e RecordLinkVisitEvent) Name() string {
	return "RecordLinkVisitEvent"
}

func NewLinkAccessedEvent(recordInfo valobj.ShortLinkStatsRecordVo) RecordLinkVisitEvent {
	return RecordLinkVisitEvent{
		CommonEvent: base_event.NewCommonEvent(),
		RecordInfo:  recordInfo,
	}
}