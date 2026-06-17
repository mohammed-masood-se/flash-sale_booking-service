package ports

import (
	"booking-service/internal/core/domain"
	"context"
)

type UserService interface {
	RegisterUser(ctx context.Context, email, password string) (insertedID string, err error)
	VerifyUser(ctx context.Context, email, code string) (insertedID string, err error)
}

type UserRepository interface {
	InsertUser(ctx context.Context, user domain.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (user *domain.User, err error)

	InsertRegistration(ctx context.Context, reg domain.Registration) (insertedID string, err error)
	GetRegistrationByEmail(ctx context.Context, email string) (reg *domain.Registration, err error)
}
