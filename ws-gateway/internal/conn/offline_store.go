package conn

import (
	"context"
	"encoding/json"
	"fmt"
	"my-IMSystem/ws-gateway/internal/model"

	"github.com/redis/go-redis/v9"
)

type RedisOfflineMsgStore struct {
	rdb *redis.Client
}

func NewRedisOfflineMsgStore(rdb *redis.Client) *RedisOfflineMsgStore {
	return &RedisOfflineMsgStore{rdb: rdb}
}

func (s *RedisOfflineMsgStore) Save(userId int64, msg model.Message) error {
	key := fmt.Sprintf("offline:msg:%d", userId)
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.rdb.RPush(context.Background(), key, data).Err()
}

func (s *RedisOfflineMsgStore) LoadAndDelete( userId int64) ([]model.Message, error) {
	key := fmt.Sprintf("offline:msg:%d", userId)
	msgs, err := s.rdb.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var result []model.Message
	for _, raw := range msgs {
		var m model.Message
		if err := json.Unmarshal([]byte(raw), &m); err == nil {
			result = append(result, m)
		}
	}

	_ = s.rdb.Del(context.Background(), key).Err()
	return result, nil
}
