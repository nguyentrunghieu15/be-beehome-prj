package datasource

import (
	"context"
	"time"

	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
)

type ISessionStorage interface {
	SaveSession(string, string, SessionKeyConfig) error
	GetSession(string) (string, error)
}

type SessionKeyConfig struct {
	ExpireTime time.Duration
}

type SessionStorage struct {
	storeRedis *database.RedisDb
}

func (s *SessionStorage) SaveSession(key, value string, config SessionKeyConfig) error {
	// Implement to save session data to Redis with expiration time
	if result := s.storeRedis.SetEx(context.Background(), key, value, config.ExpireTime); result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (s *SessionStorage) GetSession(key string) (string, error) {
	// Implement  method `Get(key string) (string, error)`:
	result := s.storeRedis.Get(context.Background(), key)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Val(), nil
}

func NewSessionStorage(redis *database.RedisDb) *SessionStorage {
	return &SessionStorage{storeRedis: redis}
}
