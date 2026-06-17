package mongorepo

import (
	"booking-service/internal/adapters/mongorepo/schema"
	"booking-service/internal/core/domain"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	usersCollection        *mongo.Collection
	registrationCollection *mongo.Collection
}

func NewUserRepository(ctx context.Context, usersCollection, registrationCollection *mongo.Collection) (*UserRepository, error) {

	usersCollectionsIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := usersCollection.Indexes().CreateOne(ctx, usersCollectionsIndex)
	if err != nil {
		return nil, fmt.Errorf("failed creating index for users collection: %w", err)
	}

	registrationsCollectionsIndex := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(120),
		},
	}
	_, err = registrationCollection.Indexes().CreateMany(ctx, registrationsCollectionsIndex)
	if err != nil {
		return nil, fmt.Errorf("failed creating indexs for registration collection: %w", err)
	}

	return &UserRepository{
		usersCollection:        usersCollection,
		registrationCollection: registrationCollection,
	}, nil
}

func (repo *UserRepository) InsertUser(ctx context.Context, user domain.User) (string, error) {
	result, err := repo.usersCollection.InsertOne(ctx, schema.User{
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", domain.ErrDuplicateEntry
		}
		return "", fmt.Errorf("failed running InsertOne: %w", err)
	}

	oid, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed converting insertedID into bson.ObjectID")
	}

	return oid.Hex(), nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user schema.User
	err := repo.usersCollection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed running FindOne: %w", err)
	}

	return user.ToDomainModel(), nil
}

func (repo *UserRepository) InsertRegistration(ctx context.Context, reg domain.Registration) (string, error) {
	result, err := repo.registrationCollection.InsertOne(ctx, schema.Registration{
		Email:    reg.Email,
		Password: reg.Password,
		Code:     reg.Code,

		CreatedAt: reg.CreatedAt,
		UpdatedAt: reg.UpdatedAt,
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", domain.ErrDuplicateEntry
		}
		return "", fmt.Errorf("failed inserting registration: %w", err)
	}

	insertedID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed converting insertedID to string: %w", err)
	}

	return insertedID.Hex(), nil
}

func (repo *UserRepository) GetRegistrationByEmail(ctx context.Context, email string) (*domain.Registration, error) {
	var registration schema.Registration
	err := repo.registrationCollection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&registration)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed running FindOne: %w", err)
	}

	return registration.ToDomainModel(), nil
}
