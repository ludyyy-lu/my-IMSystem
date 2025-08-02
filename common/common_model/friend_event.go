package common_model

type FriendEventType string

const (
	FriendRequestReceived FriendEventType = "FriendRequestReceived"
	FriendRequestAccepted FriendEventType = "FriendRequestAccepted"
	FriendRequestRejected FriendEventType = "FriendRequestRejected"
	FriendDeleted         FriendEventType = "FriendDeleted"
	UserBlocked           FriendEventType = "UserBlocked"
	UserUnblocked         FriendEventType = "UserUnblocked"
)

type FriendEvent struct {
	EventType FriendEventType `json:"event_type"`
	FromUser  int64           `json:"from_user"`
	ToUser    int64           `json:"to_user"`
	Timestamp int64           `json:"timestamp"`
	Extra     string          `json:"extra"` // 可以放验证消息等
}
