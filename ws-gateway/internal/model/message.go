package model

type Message struct {
	Type    string `json:"type"`    // 消息类型，如 chat、add_friend、heartbeat
	To      int64  `json:"to"`      // 接收方 userId
	Content string `json:"content"` // 消息内容
	From    int64  `json:"from"` // 补充发送者 ID，方便消息回显
}
