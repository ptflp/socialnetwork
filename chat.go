package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Chats struct {
	UUID      types.NullUUID  `json:"chat_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	Type      types.NullInt64 `json:"type" db:"type" ops:"update,create" orm_type:"int"`
	CreatedAt time.Time       `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
}

func (c Chats) TableName() string {
	return "chats"
}

type ChatsRepository interface {
	Update(ctx context.Context, chats Chats) error
	Find(ctx context.Context, chats Chats) (Chats, error)
	FindAll(ctx context.Context) ([]Chats, error)
	FindLimitOffset(ctx context.Context, limit, offset uint64) ([]Chats, error)
	CreateChats(ctx context.Context, chats Chats) error
}
