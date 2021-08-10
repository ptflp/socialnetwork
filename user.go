package infoblog

import (
	"context"
	"time"
)

type User struct {
	ID             int64       `json:"-" db:"id"`
	UUID           string      `json:"user_id" db:"uuid" ops:"create"`
	Phone          NullString  `json:"phone" db:"phone" ops:"update,create"`
	Email          NullString  `json:"email" db:"email" ops:"update,create"`
	Avatar         NullString  `json:"profile_image" db:"avatar" ops:"update"`
	Password       NullString  `json:"password,omitempty" db:"password" ops:"create"`
	Active         NullBool    `json:"active" db:"active" ops:"create,update"`
	Name           NullString  `json:"name" db:"name" ops:"update,create"`
	SecondName     NullString  `json:"second_name" db:"second_name" ops:"update,create"`
	EmailVerified  NullBool    `json:"email_verified" db:"email_verified"`
	Description    NullString  `json:"description" db:"description" ops:"update,create"`
	NickName       NullString  `json:"nickname" db:"nickname" ops:"update,create"`
	ShowSubs       NullBool    `json:"show_subs" db:"show_subs" ops:"update,create"`
	Cost           NullFloat64 `json:"cost" db:"cost" ops:"update,create"`
	Trial          NullBool    `json:"trial" db:"trial" ops:"update,create"`
	NotifyEmail    NullBool    `json:"notify_email" db:"notify_email" ops:"update,create"`
	NotifyTelegram NullBool    `json:"notify_telegram" db:"notify_telegram" ops:"update,create"`
	NotifyPush     NullBool    `json:"notify_push" db:"notify_push" ops:"update,create"`
	Language       NullInt64   `json:"language" db:"language" ops:"update,create"`
	FacebookID     NullInt64   `json:"facebook_id" db:"facebook_id" ops:"update,create"`
	GoogleID       NullString  `json:"google_id" db:"google_id" ops:"update,create"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
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
