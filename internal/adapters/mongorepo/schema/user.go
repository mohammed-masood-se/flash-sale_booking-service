package schema

import (
	"booking-service/internal/core/domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	Email    string        `bson:"email"`
	Password string        `bson:"password"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func (mongoUser *User) ToDomainModel() *domain.User {
	return &domain.User{
		ID:       mongoUser.ID.Hex(),
		Email:    mongoUser.Email,
		Password: mongoUser.Password,

		CreatedAt: mongoUser.CreatedAt,
		UpdatedAt: mongoUser.UpdatedAt,
	}
}
