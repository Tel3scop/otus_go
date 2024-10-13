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
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	DateTime     time.Time     `db:"datetime"`
	Duration     time.Duration `db:"duration"`
	Description  string        `db:"description"`
	UserID       string        `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}
