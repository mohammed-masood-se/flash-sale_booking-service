package domain

import (
	"errors"
	"time"
)

type User struct {
	ID         string
	Email      string
	Password   string
	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrUserAlreadyRegistered = errors.New("user has already registered")
var ErrInvalidCredentails = errors.New("invalid credentials")
