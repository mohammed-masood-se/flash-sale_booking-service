package mongorepo

import (
	"booking-service/internal/core/ports"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoTransactionManager struct {
	client *mongo.Client
}

func NewMongoTransactionManager(client *mongo.Client) *MongoTransactionManager {
	return &MongoTransactionManager{
		client: client,
	}
}

func (tm *MongoTransactionManager) Run(ctx context.Context, fn ports.TxFunc) error {
	session, err := tm.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed starting new transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(txctx context.Context) (any, error) {
		return nil, fn(txctx)
	})

	return err
}
