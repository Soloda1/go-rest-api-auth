package auth

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"go-rest-api-auth/internal/database"
	"time"

	"github.com/google/uuid"
)

type SessionManager struct {
	cacheClient        *database.CacheClient
	Ttl                time.Duration
	ErrSessionNotFound error
}

func NewSessionManager(cacheClient *database.CacheClient, ttl time.Duration) *SessionManager {
	return &SessionManager{
		cacheClient:        cacheClient,
		Ttl:                ttl,
		ErrSessionNotFound: errors.New("session not found"),
	}
}

func (sm *SessionManager) CreateSession(userID string) (string, error) {
	sessionID := uuid.New().String()
	err := sm.cacheClient.Cache.Set(sm.cacheClient.Ctx, sessionID, userID, sm.Ttl).Err()
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (sm *SessionManager) GetUserIdBySession(sessionID string) (string, error) {
	userID, err := sm.cacheClient.Cache.Get(sm.cacheClient.Ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		return "", sm.ErrSessionNotFound
	} else if err != nil {
		return "", err
	}
	return userID, nil
}

func (sm *SessionManager) GetSessionByUserID(userID string) (string, error) {
	keys, err := sm.cacheClient.Cache.Keys(sm.cacheClient.Ctx, "*").Result()
	if err != nil {
		return "", err
	}

	for _, key := range keys {
		storedUserID, err := sm.cacheClient.Cache.Get(sm.cacheClient.Ctx, key).Result()
		if err != nil {
			return "", err
		}

		if userID == storedUserID {
			return key, nil
		}
	}

	return "", sm.ErrSessionNotFound
}

func (sm *SessionManager) DeleteSession(sessionID string) error {
	err := sm.cacheClient.Cache.Del(sm.cacheClient.Ctx, sessionID).Err()
	return err
}
