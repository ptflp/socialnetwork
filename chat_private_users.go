package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type ChatPrivateUser struct {
	UserUUID   types.NullUUID `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null" orm_index:"index"`
	ToUserUUID types.NullUUID `json:"to_user_id" db:"to_user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null" orm_index:"index"`
	ChatUUID   types.NullUUID `json:"chat_id" db:"chat_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null" orm_index:"index"`
	Active     types.NullBool `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
}

func (c ChatPrivateUser) OnCreate() string {
	return ""
}

func (c ChatPrivateUser) TableName() string {
	return "chat_private_users"
}

type ChatPrivateUsersRepository interface {
	Create(ctx context.Context, chatParticipant ChatPrivateUser) error
	Find(ctx context.Context, chatParticipant ChatPrivateUser) (ChatPrivateUser, error)
	Update(ctx context.Context, chatParticipant ChatPrivateUser) error
	Delete(ctx context.Context, chatParticipant ChatPrivateUser) error
	List(ctx context.Context, limit, offset uint64) ([]ChatPrivateUser, error)
	Listx(ctx context.Context, condition Condition) ([]ChatPrivateUser, error)
}
