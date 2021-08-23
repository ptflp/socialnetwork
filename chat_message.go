package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type ChatMessages struct {
	UUID      types.NullUUID `json:"message_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	ChatUUID  types.NullUUID `json:"chat_id" db:"chat_uuid" ops:"create" orm_type:"binary(16)" orm_default:"null"`
	UserUUID  types.NullUUID `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"null"`
	Active    types.NullBool `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"not null"`
	CreatedAt time.Time      `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
	Message   string         `json:"message" db:"message" ops:"update,create" orm_type:"varchar(233)" orm_default:"not null"`
}

func (c ChatMessages) TableName() string {
	return "chatMessages"
}

type ChatMessagesRepository interface {
	Create(ctx context.Context, chatMessages ChatMessages) error
	Find(ctx context.Context, chatMessages ChatMessages) (ChatMessages, error)
	Update(ctx context.Context, chatMessages ChatMessages) error
	Delete(ctx context.Context, chatMessages ChatMessages) error
	List(ctx context.Context, limit, offset uint64) ([]ChatMessages, error)
}
