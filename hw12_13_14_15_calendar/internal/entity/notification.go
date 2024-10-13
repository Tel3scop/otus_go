package entity

import "time"

// Notification структура уведомления.
type Notification struct {
	EventID       string
	EventTitle    string
	EventDateTime time.Time
	UserID        string
}
