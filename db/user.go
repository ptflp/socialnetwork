package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	updateUser                = "UPDATE users SET email = ?, phone = ?, name = ?, second_name = ? WHERE id = ?;"
	setPassword               = "UPDATE users SET password = ? WHERE id = ?"
	findUserByPhone           = "SELECT id, email, phone FROM users WHERE phone = ?"
	findUserByEmail           = "SELECT id, email, phone, password FROM users WHERE email = ?"
	createUserByPhone         = "INSERT INTO users (phone) VALUES (?)"
	createUserByEmailPassword = "INSERT INTO users (email, password) VALUES (?, ?)"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) infoblog.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (infoblog.User, error) {

	var (
		id       sql.NullInt64
		phone    sql.NullString
		password sql.NullString
	)

	if err := u.db.QueryRowContext(ctx, findUserByEmail, email).Scan(&id, &email, &email, &password); err != nil {
		return infoblog.User{}, err
	}

	return infoblog.User{
		ID:    id.Int64,
		Email: email,
		Phone: phone.String,
	}, nil

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

func (u *userRepository) CreateUserByEmailPassword(ctx context.Context, email, passHash string) error {
	_, err := u.db.MustExecContext(ctx, createUserByEmailPassword, email, passHash).RowsAffected()

	return err
}

func (u *userRepository) Update(ctx context.Context, user infoblog.User) error {
	_, err := u.db.MustExecContext(ctx, updateUser, user.Email, user.Phone, user.Name, user.SecondName, user.ID).RowsAffected()

	return err
}

func (u *userRepository) SetPassword(ctx context.Context, user infoblog.User) error {
	_, err := u.db.MustExecContext(ctx, setPassword, user.Password, user.ID).RowsAffected()

	return err
}
