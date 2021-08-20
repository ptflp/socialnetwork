package infoblog

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type ChatMessages struct {
	UUID      types.NullUUID   `json:"uuid" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	ChatUUID  types.NullUUID   `json:"chat_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	UserUUID  types.NullUUID   `json:"user_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	CreatedAt time.Time        `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
	DeletedAt sql.NullTime     `json:"deleted_at" db:"deleted_at" orm_type:"timestamp" orm_default:"null" orm_index:"index"`
	Message   types.NullString `json:"message" db:"message" ops:"update,create" orm_type:"varchar(233)"  orm_default:"not null"`
}

func (u ChatMessages) TableName() string {
	return "chatMessages"
}

type ChatMessagesRepository interface {
	Update(ctx context.Context, chatMessages ChatMessages) error
	Find(ctx context.Context, chatMessages ChatMessages) (ChatMessages, error)
	FindAll(ctx context.Context) ([]ChatMessages, error)
	FindLimitOffset(ctx context.Context, limit, offset uint64) ([]ChatMessages, error)
	CreateChatMessages(ctx context.Context, chatMessages ChatMessages) error
}
