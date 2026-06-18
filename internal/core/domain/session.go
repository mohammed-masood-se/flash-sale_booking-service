package domain

import "time"

type Session struct {
	SessionID string
	Email     string

	CreatedAt time.Time
}
