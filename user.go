package infoblog

import "context"

type User struct {
	ID         int64  `json:"id" db:"id"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password"`
	Active     string `json:"active" db:"active"`
	Name       string `json:"name" db:"name"`
	SecondName string `json:"second_name" db:"second_name"`
}

type UserRepository interface {
	Update(ctx context.Context, user User) error
	SetPassword(ctx context.Context, user User) error

	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)

	CreateUserByPhone(ctx context.Context, phone string) error
	CreateUserByEmailPassword(ctx context.Context, email, passHash string) error
}
