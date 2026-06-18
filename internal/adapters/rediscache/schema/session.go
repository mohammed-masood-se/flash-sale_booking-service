package schema

import "time"

const SESSION_PREFIX = "session"

type Session struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
