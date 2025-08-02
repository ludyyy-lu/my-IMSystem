package common_model

type FriendEventType string

const (
	FriendRequestReceived FriendEventType = "FriendRequestReceived" // 好友请求接收
	FriendRequestAccepted FriendEventType = "FriendRequestAccepted" // 好友请求接受
	FriendRequestRejected FriendEventType = "FriendRequestRejected" // 好友请求拒绝
	FriendDeleted         FriendEventType = "FriendDeleted"         // 好友删除
	UserBlocked           FriendEventType = "UserBlocked"           // 用户被屏蔽
	UserUnblocked         FriendEventType = "UserUnblocked"         // 用户解除屏蔽
)

type FriendEvent struct {
	EventType FriendEventType `json:"event_type"`
	FromUser  int64           `json:"from_user"`
	ToUser    int64           `json:"to_user"`
	Timestamp int64           `json:"timestamp"`
	Extra     string          `json:"extra"` // 可以放验证消息等
}
