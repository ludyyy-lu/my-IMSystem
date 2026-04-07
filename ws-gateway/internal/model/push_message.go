package model

type PushMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

const (
	PushTypeChatMessage    = "chat_message"
	PushTypeFriendEvent    = "friend_event"
	PushTypeOfflineMessage = "offline_message"
)
