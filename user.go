package infoblog

import (
	"context"
	"time"
)

type User struct {
	ID             int64       `json:"-" db:"id"`
	UUID           string      `json:"uuid" db:"uuid"`
	Phone          NullString  `json:"phone" db:"phone" ops:"update"`
	Email          NullString  `json:"email" db:"email" ops:"update"`
	Password       NullString  `json:"password,omitempty" db:"password"`
	Active         NullBool    `json:"active" db:"active"`
	Name           NullString  `json:"name" db:"name" ops:"update"`
	SecondName     NullString  `json:"second_name" db:"second_name" ops:"update"`
	EmailVerified  NullBool    `json:"email_verified" db:"email_verified"`
	Description    NullString  `json:"description" db:"description" ops:"update"`
	NickName       NullString  `json:"nickname" db:"nickname" ops:"update"`
	ShowSubs       NullBool    `json:"show_subs" db:"show_subs" ops:"update"`
	Cost           NullFloat64 `json:"cost" db:"cost" ops:"update"`
	Trial          NullBool    `json:"trial" db:"trial" ops:"update"`
	NotifyEmail    NullBool    `json:"notify_email" db:"notify_email" ops:"update"`
	NotifyTelegram NullBool    `json:"notify_telegram" db:"trial" ops:"update"`
	NotifyPush     NullBool    `json:"notify_push" db:"notify_push" ops:"update"`
	Language       NullInt64   `json:"language" db:"language" ops:"update"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
}

type UserRepository interface {
	Update(ctx context.Context, user User) error
	SetPassword(ctx context.Context, user User) error

	Find(ctx context.Context, user User) (User, error)
	FindByPhone(ctx context.Context, user User) (User, error)
	FindByEmail(ctx context.Context, user User) (User, error)

	CreateUserByPhone(ctx context.Context, user User) error
	CreateUserByEmailPassword(ctx context.Context, user User) error
}
