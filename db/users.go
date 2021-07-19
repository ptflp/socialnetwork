package db

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	updateUser  = "UPDATE users SET phone = ?, email = ?, name = ?, second_name = ?, email_verified = ? WHERE uuid = ?;"
	setPassword = "UPDATE users SET password = ? WHERE uuid = ?"

	findUserByPhone = "SELECT id, email, phone FROM users WHERE phone = ?"
	findUserByEmail = "SELECT id, email, phone, password, email_verified FROM users WHERE email = ?"

	createUserByPhone         = "INSERT INTO users (uuid, active, phone) VALUES (?, ?, ?)"
	createUserByEmailPassword = "INSERT INTO users (uuid, email, password, active, email_verified) VALUES (?, ?, ?, 1, 1)"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) infoblog.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) FindByEmail(ctx context.Context, user infoblog.User) (infoblog.User, error) {

	fields, err := infoblog.GetFields("users")
	if err != nil {
		return infoblog.User{}, err
	}

	query, args, err := sq.Select(fields...).From("users").Where(sq.Eq{"email": user.Email}).ToSql()
	if err != nil {
		return infoblog.User{}, err
	}

	if err := u.db.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		return infoblog.User{}, err
	}

	return user, nil
}

func (u *userRepository) FindByPhone(ctx context.Context, user infoblog.User) (infoblog.User, error) {

	fields, err := infoblog.GetFields("users")
	if err != nil {
		return infoblog.User{}, err
	}

	query, args, err := sq.Select(fields...).From("users").Where(sq.Eq{"phone": user.Phone}).ToSql()
	if err != nil {
		return infoblog.User{}, err
	}

	if err = u.db.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		return infoblog.User{}, err
	}

	return user, nil
}

func (u *userRepository) CreateUserByPhone(ctx context.Context, user infoblog.User) error {
	if !user.Phone.Valid {
		return fmt.Errorf("bad phone number %s", user.Phone.String)
	}
	_, err := u.db.MustExecContext(ctx, createUserByPhone, user.UUID, true, user.Phone).RowsAffected()

	return err
}

func (u *userRepository) CreateUserByEmailPassword(ctx context.Context, user infoblog.User) error {
	err := u.db.QueryRowContext(ctx, createUserByEmailPassword, user.UUID, user.Email, user.Password).Err()

	return err
}

func (u *userRepository) Update(ctx context.Context, user infoblog.User) error {
	if len(user.UUID) != 40 {
		return errors.New("wrong user uuid on update")
	}
	updateFields, err := infoblog.GetUpdateFields("users")
	if err != nil {
		return err
	}
	updateFieldsPointers := infoblog.GetFieldsPointers(&user, "update")

	queryRaw := sq.Update("users").Where(sq.Eq{"uuid": user.UUID})
	for i := range updateFields {
		queryRaw = queryRaw.Set(updateFields[i], updateFieldsPointers[i])
	}

	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}
	res, err := u.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()

	return err
}

func (u *userRepository) SetPassword(ctx context.Context, user infoblog.User) error {
	_, err := u.db.MustExecContext(ctx, setPassword, user.Password, user.UUID).RowsAffected()

	return err
}

func (u *userRepository) Find(ctx context.Context, user infoblog.User) (infoblog.User, error) {
	fields, err := infoblog.GetFields("users")
	if err != nil {
		return infoblog.User{}, err
	}

	query, args, err := sq.Select(fields...).From("users").Where(sq.Eq{"uuid": user.UUID}).ToSql()
	if err != nil {
		return infoblog.User{}, err
	}

	if err := u.db.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		return infoblog.User{}, err
	}

	return user, nil
}
