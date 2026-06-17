package schema

import (
	"booking-service/internal/core/domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Registration struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	Email    string        `bson:"email"`
	Password string        `bson:"password"`
	Code     string        `bson:"code"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func (registration *Registration) ToDomainModel() *domain.Registration {
	return &domain.Registration{
		ID:       registration.ID.Hex(),
		Email:    registration.Email,
		Password: registration.Password,
		Code:     registration.Code,

		CreatedAt: registration.CreatedAt,
		UpdatedAt: registration.UpdatedAt,
	}
}
