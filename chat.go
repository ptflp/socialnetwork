package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Chat struct {
	UUID            types.NullUUID  `json:"chat_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	Type            types.NullInt64 `json:"type" db:"type" ops:"update,create" orm_type:"int" `
	Active          types.NullBool  `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	UserUUID        types.NullUUID  `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null" orm_index:"index"`
	LastMessageUUID types.NullUUID  `json:"last_message" db:"last_message" ops:"create,update" orm_type:"binary(16)" orm_default:"null" orm_index:"index"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
}

func (c Chat) OnCreate() string {
	return ""
}

func (c Chat) TableName() string {
	return "chats"
}

type ChatRepository interface {
	Create(ctx context.Context, chat Chat) error
	Find(ctx context.Context, chat Chat) (Chat, error)
	Update(ctx context.Context, chat Chat) error
	Delete(ctx context.Context, chat Chat) error
	List(ctx context.Context, limit, offset uint64) ([]Chat, error)
	Listx(ctx context.Context, condition Condition) ([]Chat, error)
}
