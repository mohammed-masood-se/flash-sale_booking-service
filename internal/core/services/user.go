package services

import (
	"booking-service/internal/core/domain"
	"booking-service/internal/core/ports"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

type UserService struct {
	txmanager      ports.TxManager
	userRepository ports.UserRepository
	userCache      ports.UserCache
}

func NewUserService(txmanager ports.TxManager, userRepository ports.UserRepository, userCache ports.UserCache) *UserService {
	return &UserService{
		txmanager:      txmanager,
		userRepository: userRepository,
		userCache:      userCache,
	}
}

func (service *UserService) RegisterUser(ctx context.Context, email, password string) (string, error) {

	currentTime := time.Now().UTC()

	registration := domain.Registration{
		Email:     email,
		Password:  password,
		Code:      "1011",
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	var insertedID string

	err := service.txmanager.Run(ctx, func(ctx context.Context) error {
		_, err := service.userRepository.GetUserByEmail(ctx, email)
		if err == nil {
			return domain.NewClientError(domain.ErrUserAlreadyRegistered)
		}
		if !errors.Is(err, domain.ErrNotFound) {
			return domain.NewServiceError(fmt.Errorf("unexpected error when getting user by email: %w", err))
		}

		insertedID, err = service.userRepository.InsertRegistration(ctx, registration)
		if err != nil {
			if errors.Is(err, domain.ErrDuplicateEntry) {
				return domain.NewClientError(domain.ErrRegistrationAlreadyInProgress)
			}
			return domain.NewServiceError(fmt.Errorf("unexpected error when inserting registration: %w", err))
		}

		return nil
	})

	return insertedID, err
}

func (service *UserService) VerifyUser(ctx context.Context, email, code string) (string, error) {

	var insertedID string

	err := service.txmanager.Run(ctx, func(ctx context.Context) error {

		registration, err := service.userRepository.GetRegistrationByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return domain.NewClientError(domain.ErrNotFound)
			}
			return domain.NewServiceError(fmt.Errorf("unexpected error when getting registration by email: %w", err))
		}

		if registration.Code != code {
			return domain.NewClientError(domain.ErrInvalidVerificationCode)
		}

		insertedID, err = service.userRepository.InsertUser(ctx, domain.User{
			Email:    registration.Email,
			Password: registration.Password,
		})

		if err != nil {
			if errors.Is(err, domain.ErrDuplicateEntry) {
				return domain.NewClientError(domain.ErrDuplicateEntry)
			}
			return domain.NewServiceError(fmt.Errorf("unexpected error when inserting user: %w", err))
		}

		return nil
	})

	return insertedID, err
}

func (service *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {

	user, err := service.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.NewClientError(domain.ErrNotFound)
		}
		return "", domain.NewServiceError(fmt.Errorf("unexpected error while getting user by email: %w", err))
	}

	if user.Password != password {
		return "", domain.NewClientError(domain.ErrInvalidCredentails)
	}

	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", domain.NewServiceError(fmt.Errorf("failed creating sessionID: %w", err))
	}

	sessionID := base64.RawURLEncoding.EncodeToString(bytes)

	currentTime := time.Now().UTC()
	err = service.userCache.InsertSession(ctx, domain.Session{
		SessionID: sessionID,
		Email:     email,
		CreatedAt: currentTime,
	})
	if err != nil {
		return "", domain.NewServiceError(fmt.Errorf("unexpected error while inserting user session: %w", err))
	}

	return sessionID, nil
}
