package logic

import (
	"log"
	"my-IMSystem/ws-gateway/internal/model"
	"my-IMSystem/ws-gateway/internal/svc"
)

// MessageHandler 定义每种消息的处理函数签名
type MessageHandler func(svcCtx *svc.ServiceContext, fromUserId int64, msg model.Message)

// 消息类型路由表
var messageRouter = map[string]MessageHandler{
	"chat": handleChatMessage,
	// 未来可以加： "add_friend": handleAddFriendMessage,
}

func RouteMessage(svcCtx *svc.ServiceContext, fromUserId int64, msg model.Message) {
	handler, ok := messageRouter[msg.Type]
	if !ok {
		log.Printf("Unknown message type: %s", msg.Type)
		return
	}
	handler(svcCtx, fromUserId, msg)
}
