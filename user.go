package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type User struct {
	UUID           types.NullUUID    `json:"user_id" db:"uuid" ops:"create"`
	Phone          types.NullString  `json:"phone" db:"phone" ops:"update,create"`
	Email          types.NullString  `json:"email" db:"email" ops:"update,create"`
	Avatar         types.NullString  `json:"profile_image" db:"avatar" ops:"update"`
	Password       types.NullString  `json:"password,omitempty" db:"password" ops:"create"`
	Active         types.NullBool    `json:"active" db:"active" ops:"create,update"`
	Name           types.NullString  `json:"name" db:"name" ops:"update,create"`
	SecondName     types.NullString  `json:"second_name" db:"second_name" ops:"update,create"`
	EmailVerified  types.NullBool    `json:"email_verified" db:"email_verified"`
	Description    types.NullString  `json:"description" db:"description" ops:"update,create"`
	NickName       types.NullString  `json:"nickname" db:"nickname" ops:"update,create"`
	ShowSubs       types.NullBool    `json:"show_subs" db:"show_subs" ops:"update,create"`
	Cost           types.NullFloat64 `json:"cost" db:"cost" ops:"update,create"`
	Trial          types.NullBool    `json:"trial" db:"trial" ops:"update,create"`
	NotifyEmail    types.NullBool    `json:"notify_email" db:"notify_email" ops:"update,create"`
	NotifyTelegram types.NullBool    `json:"notify_telegram" db:"notify_telegram" ops:"update,create"`
	NotifyPush     types.NullBool    `json:"notify_push" db:"notify_push" ops:"update,create"`
	Language       types.NullInt64   `json:"language" db:"language" ops:"update,create"`
	FacebookID     types.NullInt64   `json:"facebook_id" db:"facebook_id" ops:"update,create"`
	GoogleID       types.NullString  `json:"google_id" db:"google_id" ops:"update,create"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`
}

type UserRepository interface {
	Update(ctx context.Context, user User) error
	SetPassword(ctx context.Context, user User) error

	Find(ctx context.Context, user User) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindLimitOffset(ctx context.Context, limit, offset uint64) ([]User, error)
	FindByPhone(ctx context.Context, user User) (User, error)
	FindByEmail(ctx context.Context, user User) (User, error)
	FindByNickname(ctx context.Context, user User) (User, error)
	FindLikeNickname(ctx context.Context, nickname string) ([]User, error)
	FindByFacebook(ctx context.Context, user User) (User, error)
	FindByGoogle(ctx context.Context, user User) (User, error)

	CreateUser(ctx context.Context, user User) error
	CreateUserByEmailPassword(ctx context.Context, user User) error
}
