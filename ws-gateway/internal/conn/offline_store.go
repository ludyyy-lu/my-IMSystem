package conn

import (
	"my-IMSystem/ws-gateway/internal/model"
	"sync"
)

// 离线消息存储结构（模拟 Redis）
type OfflineMsgManager struct {
	data map[int64][]model.Message // key: 用户ID，value: 消息列表
	lock sync.RWMutex
}

// 创建实例
func NewOfflineMsgManager() *OfflineMsgManager {
	return &OfflineMsgManager{
		data: make(map[int64][]model.Message),
	}
}

// 添加离线消息
func (m *OfflineMsgManager) Add(userId int64, msg model.Message) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[userId] = append(m.data[userId], msg)
}

// 获取并删除离线消息
func (m *OfflineMsgManager) Take(userId int64) []model.Message {
	m.lock.Lock()
	defer m.lock.Unlock()
	msgs := m.data[userId]
	delete(m.data, userId) // 清除
	return msgs
}
