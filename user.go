package infoblog

import "context"

type User struct {
	ID    int64  `json:"id" db:"id"`
	Phone string `json:"phone" db:"phone"`
	Email string `json:"email" db:"email"`
}

type UserRepository interface {
	FindByPhone(ctx context.Context, phone string) (User, error)
	CreateUserByPhone(ctx context.Context, phone string) error
}
