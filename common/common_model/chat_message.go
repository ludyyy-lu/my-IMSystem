package common_model

// Kafka 投递的数据格式
type ChatMessage struct {
	MessageId   string `json:"message_id"`
	FromUserId  int64  `json:"from_user_id"`
	ToUserId    int64  `json:"to_user_id"`
	Content     string `json:"content"`
	Timestamp   int64  `json:"timestamp"`
}
