package logic

import (
	"context"
	"my-IMSystem/chat-service/internal/model"

	"gorm.io/gorm"
)

// 把消息存入 MySQL
type SaveChatMessageLogic struct {
	ctx context.Context
	db  *gorm.DB
}

func NewSaveChatMessageLogic(ctx context.Context, db *gorm.DB) *SaveChatMessageLogic {
	return &SaveChatMessageLogic{
		ctx: ctx,
		db:  db,
	}
}

func (l *SaveChatMessageLogic) Save(msg *model.Message) error {
	return l.db.WithContext(l.ctx).Create(msg).Error
}
