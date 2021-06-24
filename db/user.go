package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/ptflp/infoblog-server"
)

const (
	findUserByPhone   = "SELECT id, email FROM users WHERE phone = ?"
	createUserByPhone = "INSERT INTO users (phone) VALUES (?)"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) infoblog.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) FindByPhone(ctx context.Context, phone string) (infoblog.User, error) {

	var (
		id    sql.NullInt64
		email sql.NullString
	)

	if err := u.db.QueryRowContext(ctx, findUserByPhone, phone).Scan(&id, &email, &phone); err != nil {
		return infoblog.User{}, err
	}

	return infoblog.User{
		ID:    id.Int64,
		Email: email.String,
		Phone: phone,
	}, nil
}

func (u *userRepository) CreateUserByPhone(ctx context.Context, phone string) error {
	_, err := u.db.MustExecContext(ctx, createUserByPhone, phone).RowsAffected()

	return err
}
