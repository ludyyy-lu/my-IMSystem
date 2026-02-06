package model

type PushMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
