package transport

import (
	"encoding/json"

	"my-IMSystem/ws-gateway/internal/model"
)

func ParseMessage(payload []byte) (model.WsMessage, error) {
	var msg model.WsMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return model.WsMessage{}, err
	}
	return msg, nil
}
