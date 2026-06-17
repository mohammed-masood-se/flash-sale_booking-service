package domain

import (
	"errors"
	"time"
)

type Registration struct {
	ID       string
	Email    string
	Password string
	Code     string

	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrRegistrationAlreadyInProgress = errors.New("registration is already in progress")
var ErrInvalidVerificationCode = errors.New("invalid verification code")
