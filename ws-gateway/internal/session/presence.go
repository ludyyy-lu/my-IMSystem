// Package session – presence store tracks each user's online/offline status
// in Redis so that any gateway node (or downstream service) can query whether
// a user is currently connected, enabling multi-node deployments to share
// presence information without direct inter-node communication.
package session

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// presenceKeyTTL is the Redis key lifetime for an "online" entry.
	// It acts as a safety net: if a gateway crashes without calling SetOffline,
	// the key expires automatically and the user is no longer considered online.
	presenceKeyTTL = 2 * time.Hour
)

// PresenceStore tracks users' online/offline status.
// Implementations must be safe for concurrent use.
type PresenceStore interface {
	// SetOnline marks the user as online and refreshes the entry TTL.
	SetOnline(ctx context.Context, userID int64) error

	// SetOffline removes the user's online entry immediately.
	SetOffline(ctx context.Context, userID int64) error

	// IsOnline reports whether the user currently has an active session on
	// any gateway node that has called SetOnline.
	IsOnline(ctx context.Context, userID int64) (bool, error)
}

// RedisPresenceStore implements PresenceStore using a single Redis key per user.
// Key schema: presence:online:{userID}
type RedisPresenceStore struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewRedisPresenceStore creates a RedisPresenceStore backed by the given client.
// The default TTL (presenceKeyTTL) is used for each online entry.
func NewRedisPresenceStore(rdb *redis.Client) *RedisPresenceStore {
	return &RedisPresenceStore{rdb: rdb, ttl: presenceKeyTTL}
}

func presenceKey(userID int64) string {
	return fmt.Sprintf("presence:online:%d", userID)
}

// SetOnline writes (or refreshes) the user's presence key with a TTL.
func (s *RedisPresenceStore) SetOnline(ctx context.Context, userID int64) error {
	return s.rdb.Set(ctx, presenceKey(userID), 1, s.ttl).Err()
}

// SetOffline deletes the user's presence key.
func (s *RedisPresenceStore) SetOffline(ctx context.Context, userID int64) error {
	return s.rdb.Del(ctx, presenceKey(userID)).Err()
}

// IsOnline returns true if a presence key exists for the user.
func (s *RedisPresenceStore) IsOnline(ctx context.Context, userID int64) (bool, error) {
	n, err := s.rdb.Exists(ctx, presenceKey(userID)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
