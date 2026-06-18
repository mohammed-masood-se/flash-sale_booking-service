package rediscache

import (
	"booking-service/internal/adapters/rediscache/schema"
	"booking-service/internal/core/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	redisClient *redis.Client
}

func NewUserCache(redisClient *redis.Client) *UserCache {
	return &UserCache{
		redisClient: redisClient,
	}
}

func (cache *UserCache) InsertSession(ctx context.Context, session domain.Session) error {

	objBytes, err := json.Marshal(schema.Session{
		Email:     session.Email,
		CreatedAt: session.CreatedAt,
	})

	if err != nil {
		return fmt.Errorf("failed marshalling session schema: %w", err)
	}

	err = cache.redisClient.Set(ctx, fmt.Sprintf("%v:%v", schema.SESSION_PREFIX, session.SessionID), objBytes, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed inserting session: %w", err)
	}

	return nil
}
