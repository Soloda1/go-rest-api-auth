package auth

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"go-rest-api-auth/internal/database"
	"time"

	"github.com/google/uuid"
)

type SessionManagerImplementation struct {
	cacheClient        *database.CacheClient
	Ttl                time.Duration
	ErrSessionNotFound error
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --name SessionManager --output ../../../testing/mocks
type SessionManager interface {
	CreateSession(userID string) (string, error)
	GetUserIdBySession(sessionID string) (string, error)
	GetSessionByUserID(userID string) (string, error)
	DeleteSession(sessionID string) error
	GetterTtl() time.Duration
	GetterErrSessionNotFound() error
}

func NewSessionManager(cacheClient *database.CacheClient, ttl time.Duration) SessionManager {
	return &SessionManagerImplementation{
		cacheClient:        cacheClient,
		Ttl:                ttl,
		ErrSessionNotFound: errors.New("session not found"),
	}
}

func (sm *SessionManagerImplementation) GetterTtl() time.Duration {
	return sm.Ttl
}

func (sm *SessionManagerImplementation) GetterErrSessionNotFound() error {
	return sm.ErrSessionNotFound
}

func (sm *SessionManagerImplementation) CreateSession(userID string) (string, error) {
	sessionID := uuid.New().String()
	err := sm.cacheClient.Cache.Set(sm.cacheClient.Ctx, sessionID, userID, sm.Ttl).Err()
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (sm *SessionManagerImplementation) GetUserIdBySession(sessionID string) (string, error) {
	userID, err := sm.cacheClient.Cache.Get(sm.cacheClient.Ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		return "", sm.ErrSessionNotFound
	} else if err != nil {
		return "", err
	}
	return userID, nil
}

func (sm *SessionManagerImplementation) GetSessionByUserID(userID string) (string, error) {
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

func (sm *SessionManagerImplementation) DeleteSession(sessionID string) error {
	err := sm.cacheClient.Cache.Del(sm.cacheClient.Ctx, sessionID).Err()
	return err
}
