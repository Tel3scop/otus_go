package entity

import "time"

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
