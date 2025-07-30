package model

type Message struct {
	FromUserId int64  `json:"from_user_id"`
	ToUserId   int64  `json:"to_user_id"`
	Content    string `json:"content"`
	MsgType    string `json:"msg_type"` // text, image, etc.
	Timestamp  int64  `json:"timestamp"`
}
