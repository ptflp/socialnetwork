package infoblog

import (
	"context"
	"time"
)

type User struct {
	ID            int64      `json:"id" db:"id"`
	UUID          string     `json:"uuid" db:"uuid"`
	Phone         NullString `json:"phone" db:"phone"`
	Email         NullString `json:"email" db:"email"`
	Password      NullString `json:"password,omitempty" db:"password"`
	Active        NullBool   `json:"active" db:"active"`
	Name          NullString `json:"name" db:"name"`
	SecondName    NullString `json:"second_name" db:"second_name"`
	EmailVerified NullBool   `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
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
