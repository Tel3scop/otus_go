package entity

import "time"

type PeriodType string

const (
	PeriodDay   = PeriodType("day")
	PeriodWeek  = PeriodType("week")
	PeriodMonth = PeriodType("month")
)

// Event структура события.
type Event struct {
	ID           string
	Title        string
	DateTime     time.Time
	Duration     time.Duration
	Description  string
	UserID       string
	NotifyBefore time.Duration
}
