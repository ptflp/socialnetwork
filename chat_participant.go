package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type ChatParticipant struct {
	ChatUUID types.NullUUID  `json:"chat_id" db:"chat_uuid" ops:"create" orm_type:"binary(16)" orm_default:"null"`
	UserUUID types.NullUUID  `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"null"`
	Type     types.NullInt64 `json:"type" db:"type" ops:"update,create" orm_type:"int" orm_default:"null"`
	Active   types.NullBool  `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	JoinedAt time.Time       `json:"joined_at" db:"joined_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
}

func (c ChatParticipant) OnCreate() string {
	return ""
}

func (c ChatParticipant) TableName() string {
	return "chat_participants"
}

type ChatParticipantRepository interface {
	Create(ctx context.Context, chatParticipant ChatParticipant) error
	Find(ctx context.Context, chatParticipant ChatParticipant) (ChatParticipant, error)
	Update(ctx context.Context, chatParticipant ChatParticipant) error
	Delete(ctx context.Context, chatParticipant ChatParticipant) error
	List(ctx context.Context, limit, offset uint64) ([]ChatParticipant, error)
	Listx(ctx context.Context, condition Condition) ([]ChatParticipant, error)
}
