package types

import "log/slog"

type PageReq struct {
	Current int `json:"current"`
	Size    int `json:"size"`
}

func (p PageReq) Limit() int {
	return p.Size
}

func (p PageReq) Offset() int {
	if p.Current <= 0 {
		return 0
	}
	return (p.Current - 1) * p.Size
}

type PageResp[T any] struct {
	Total   int64 `json:"total"`
	Current int   `json:"current"`
	Size    int   `json:"size"`
	Records []T   `json:"records"`
}

func NewEmptyPageResp[T any]() *PageResp[T] {
	return &PageResp[T]{
		Records: make([]T, 0),
	}
}

func (r PageResp[T]) WithTotal(total int64) PageResp[T] {
	r.Total = total
	return r
}

func (r PageResp[T]) WithCurrent(current int) PageResp[T] {
	r.Current = current
	return r
}

func (r PageResp[T]) WithSize(size int) PageResp[T] {
	r.Size = size
	return r
}

func (r PageResp[T]) WithRecords(records []T) PageResp[T] {
	r.Records = records
	return r
}

// ConvertRecords converts the Records field to a different type using the provided function
func ConvertRecords[S, D any](p *PageResp[S], fn func(S) (D, error)) *PageResp[D] {
	records := make([]D, 0, len(p.Records))
	for _, before := range p.Records {
		after, err := fn(before)
		if err != nil {
			slog.Error("convert records failed", "record", before, "error", err)
			continue
		}
		records = append(records, after)
	}
	return &PageResp[D]{
		Total:   p.Total,
		Current: p.Current,
		Size:    p.Size,
		Records: records,
	}
}
