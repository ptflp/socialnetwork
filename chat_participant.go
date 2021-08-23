package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type ChatParticipants struct {
	ChatUUID types.NullUUID `json:"chat_id" db:"chat_uuid" ops:"create" orm_type:"binary(16)" orm_default:"null"`
	UserUUID types.NullUUID `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"null"`
	JoinedAt time.Time      `json:"joined_at" db:"joined_at" orm_type:"timestamp" orm_default:"CURRENT_TIMESTAMP not null" orm_index:"index"`
}

func (c ChatParticipants) TableName() string {
	return "chatParticipants"
}

type ChatParticipantsRepository interface {
	Update(ctx context.Context, chatParticipants ChatParticipants) error
	Find(ctx context.Context, chatParticipants ChatParticipants) (ChatParticipants, error)
	FindAll(ctx context.Context) ([]ChatParticipants, error)
	FindLimitOffset(ctx context.Context, limit, offset uint64) ([]ChatParticipants, error)
	CreateChatParticipants(ctx context.Context, chatParticipants ChatParticipants) error
}
