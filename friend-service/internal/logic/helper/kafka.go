package helper

import (
	"encoding/json"
	"my-IMSystem/common/kafka"
	"my-IMSystem/common/common_model"
	"time"
)

func SendFriendEventToKafka(eventType common_model.FriendEventType, fromUser, toUser int64, extra string) error {
	event := common_model.FriendEvent{
		EventType: eventType,
		FromUser:  fromUser,
		ToUser:    toUser,
		Timestamp: time.Now().Unix(),
		Extra:     extra,
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return kafka.SendMessage("friend-events", data)
}
