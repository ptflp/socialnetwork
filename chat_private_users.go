package infoblog

import (
	"context"

	"gitlab.com/InfoBlogFriends/server/types"
)

type ChatPrivateUser struct {
	UserUUID   types.NullUUID `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null" orm_index:"index"`
	ToUserUUID types.NullUUID `json:"to_user_id" db:"to_user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null" orm_index:"index"`
	ChatUUID   types.NullUUID `json:"chat_id" db:"chat_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null" orm_index:"index"`
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
