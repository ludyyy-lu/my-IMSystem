// Package session provides the offline message store used to buffer messages
// for users who are not currently connected.  Messages are stored as raw bytes
// (pre-serialised JSON envelopes) so the store is agnostic to message type.
package session

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// OfflineStore is the interface for persisting and retrieving offline messages.
// Save stores a pre-serialised push-message envelope for the given user.
// LoadAndDelete atomically returns all stored envelopes and removes them from
// the store.  Implementations must be safe for concurrent use.
type OfflineStore interface {
	Save(userID int64, data []byte) error
	LoadAndDelete(userID int64) ([][]byte, error)
}

// RedisOfflineMsgStore implements OfflineStore using a Redis list per user.
// Each list element is a raw JSON-encoded push-message envelope ([]byte).
type RedisOfflineMsgStore struct {
	rdb *redis.Client
}

// NewRedisOfflineMsgStore creates a new store backed by the given Redis client.
func NewRedisOfflineMsgStore(rdb *redis.Client) *RedisOfflineMsgStore {
	return &RedisOfflineMsgStore{rdb: rdb}
}

// Save appends data to the user's offline message list.
func (s *RedisOfflineMsgStore) Save(userID int64, data []byte) error {
	key := fmt.Sprintf("offline:msg:%d", userID)
	return s.rdb.RPush(context.Background(), key, data).Err()
}

// LoadAndDelete returns all buffered envelopes and atomically removes the list.
func (s *RedisOfflineMsgStore) LoadAndDelete(userID int64) ([][]byte, error) {
	key := fmt.Sprintf("offline:msg:%d", userID)
	ctx := context.Background()

	msgs, err := s.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	result := make([][]byte, 0, len(msgs))
	for _, raw := range msgs {
		result = append(result, []byte(raw))
	}

	_ = s.rdb.Del(ctx, key).Err()
	return result, nil
}
