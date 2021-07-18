package infoblog

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID            int64          `json:"id" db:"id"`
	UUID          string         `json:"uuid" db:"uuid"`
	Phone         sql.NullString `json:"phone" db:"phone"`
	Email         sql.NullString `json:"email" db:"email"`
	Password      sql.NullString `json:"password,omitempty" db:"password"`
	Active        sql.NullBool   `json:"active" db:"active"`
	Name          sql.NullString `json:"name" db:"name"`
	SecondName    sql.NullString `json:"second_name" db:"second_name"`
	EmailVerified sql.NullBool   `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
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
