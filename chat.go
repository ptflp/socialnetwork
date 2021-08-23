package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Chat struct {
	UUID      types.NullUUID  `json:"chat_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	Type      types.NullInt64 `json:"type" db:"type" ops:"update,create" orm_type:"int"`
	CreatedAt time.Time       `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
}

func (c Chat) TableName() string {
	return "chats"
}

type ChatRepository interface {
	Create(ctx context.Context, c Chat) (int64, error)
	Update(ctx context.Context, c Chat) error
	Delete(ctx context.Context, c Chat) error

	Find(ctx context.Context, c Chat) (Chat, error)
	FindAll(ctx context.Context) ([]Chat, error)
	FindLimitOffset(ctx context.Context, limit, offset uint64) ([]Chat, error)
}
