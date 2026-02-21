package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	redisClient *redisv9.Client
}

func NewSessionRepository(redisClient *redisv9.Client) *sessionRepository {
	return &sessionRepository{redisClient: redisClient}
}

type refreshSession struct {
    JTI  string `json:"jti"`
    Role string `json:"role"`
}

type SessionRepository interface {
	SaveRefreshSession(ctx context.Context, userID, sessionID, jti, role string, ttlSeconds int64) error
	ValidateRefreshSession(ctx context.Context, userID, sessionID, jti string) (bool, error)
	DeleteRefreshSession(ctx context.Context, userID, sessionID string) error

	BlacklistAccessJTI(ctx context.Context, jti string, ttlSeconds int64) error
	IsAccessJTIBlacklisted(ctx context.Context, jti string) (bool, error)
}


func (r *sessionRepository) SaveRefreshSession(ctx context.Context, userID, sessionID, jti, role string, ttlSeconds int64) error {
    key := fmt.Sprintf("auth:refresh:%s:%s", userID, sessionID)
    value, _ := json.Marshal(refreshSession{JTI: jti, Role: role})
    return r.redisClient.Set(ctx, key, string(value), time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *sessionRepository) ValidateRefreshSession(ctx context.Context, userID, sessionID, jti string) (bool, error) {
    key := fmt.Sprintf("auth:refresh:%s:%s", userID, sessionID)
    value, err := r.redisClient.Get(ctx, key).Result()
    if err == redisv9.Nil {
        return false, nil
    }
    if err != nil {
        return false, err
    }

    var s refreshSession
    if err := json.Unmarshal([]byte(value), &s); err != nil {
        return false, err
    }

    return s.JTI == jti, nil
}

func (r *sessionRepository) DeleteRefreshSession(ctx context.Context, userID, sessionID string) error {
    key := fmt.Sprintf("auth:refresh:%s:%s", userID, sessionID)
    return r.redisClient.Del(ctx, key).Err()
}

func (r *sessionRepository) BlacklistAccessJTI(ctx context.Context, jti string, ttlSeconds int64) error {
    key := fmt.Sprintf("auth:blacklist:access:%s", jti)
    return r.redisClient.Set(ctx, key, "1", time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *sessionRepository) IsAccessJTIBlacklisted(ctx context.Context, jti string) (bool, error) {
    key := fmt.Sprintf("auth:blacklist:access:%s", jti)
    result, err := r.redisClient.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return result > 0, nil
}