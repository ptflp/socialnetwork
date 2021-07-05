package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	updateUser  = "UPDATE users SET phone = ?, email = ?, name = ?, second_name = ?, email_verified = ? WHERE id = ?;"
	setPassword = "UPDATE users SET password = ? WHERE id = ?"

	findUserByID    = "SELECT phone, email, active, name, second_name FROM users WHERE id = ?"
	findUserByPhone = "SELECT id, email, phone FROM users WHERE phone = ?"
	findUserByEmail = "SELECT id, email, phone, password, email_verified FROM users WHERE email = ?"

	createUserByPhone         = "INSERT INTO users (phone, active) VALUES (?, 1)"
	createUserByEmailPassword = "INSERT INTO users (email, password, active, email_verified) VALUES (?, ?, 1, 1)"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) infoblog.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (infoblog.User, error) {

	var (
		id            sql.NullInt64
		phone         sql.NullString
		password      sql.NullString
		emailVerified sql.NullInt64
	)

	if err := u.db.QueryRowContext(ctx, findUserByEmail, email).Scan(&id, &email, &phone, &password, &emailVerified); err != nil {
		return infoblog.User{}, err
	}

	return infoblog.User{
		ID:            id.Int64,
		Email:         email,
		Phone:         phone.String,
		EmailVerified: emailVerified.Int64,
		Password:      password.String,
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
	err := u.db.QueryRowContext(ctx, createUserByEmailPassword, email, passHash).Err()

	return err
}

func (u *userRepository) Update(ctx context.Context, user infoblog.User) error {
	if user.ID == 0 {
		return errors.New("wrong user.ID on update")
	}
	_, err := u.db.MustExecContext(ctx, updateUser, user.Phone, user.Email, user.Name, user.SecondName, user.EmailVerified, user.ID).RowsAffected()

	return err
}

func (u *userRepository) SetPassword(ctx context.Context, user infoblog.User) error {
	_, err := u.db.MustExecContext(ctx, setPassword, user.Password, user.ID).RowsAffected()

	return err
}

func (u *userRepository) Find(ctx context.Context, uid int64) (infoblog.User, error) {

	var (
		phone      sql.NullString
		email      sql.NullString
		active     sql.NullInt64
		name       sql.NullString
		secondName sql.NullString
	)

	if err := u.db.QueryRowContext(ctx, findUserByID, uid).Scan(&phone, &email, &active, &name, &secondName); err != nil {
		return infoblog.User{}, err
	}

	return infoblog.User{
		ID:         uid,
		Phone:      phone.String,
		Email:      email.String,
		Password:   "",
		Active:     active.Int64,
		Name:       name.String,
		SecondName: secondName.String,
	}, nil
}
