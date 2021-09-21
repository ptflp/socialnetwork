package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	createSubscribe = "INSERT INTO subscribes (user_uuid, subscriber_uuid, active) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE active = 1"
	deleteSubscribe = "INSERT INTO subscribes (user_uuid, subscriber_uuid, active) VALUES (?, ?, false) ON DUPLICATE KEY UPDATE active = false"
)

type subsRepository struct {
	db *sqlx.DB
	crud
}

func (sb *subsRepository) Create(ctx context.Context, sub infoblog.Subscriber) (int64, error) {
	res, err := sb.db.ExecContext(ctx, createSubscribe, sub.UserUUID, sub.SubscriberUUID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (sb *subsRepository) Delete(ctx context.Context, sub infoblog.Subscriber) error {
	_, err := sb.db.ExecContext(ctx, deleteSubscribe, sub.UserUUID, sub.SubscriberUUID)
	if err != nil {
		return err
	}

	return nil
}

func (sb *subsRepository) FindByUser(ctx context.Context, user infoblog.User) ([]infoblog.Subscriber, error) {

	panic("implement me")
}

func (sb *subsRepository) CountByUser(ctx context.Context, user infoblog.User) (int64, error) {

	var count sql.NullInt64

	query, args, err := sq.Select("COUNT(user_uuid)").From("subscribes").Where(sq.Eq{"user_uuid": user.UUID}).ToSql()
	if err != nil {
		return count.Int64, err
	}

	err = sb.db.QueryRowContext(ctx, query, args...).Scan(&count)

	return count.Int64, err
}

func (sb *subsRepository) CheckSubscribed(ctx context.Context, user infoblog.User, subscriber infoblog.User) bool {
	query, args, _ := sq.Select("active").From("subscribes").Where(sq.Eq{"subscriber_uuid": subscriber.UUID, "active": 1}).ToSql()

	res, err := sb.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false
	}

	return n > 1
}

func (sb *subsRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.Subscriber, error) {
	var subscribers []infoblog.Subscriber
	err := sb.crud.listx(ctx, &subscribers, infoblog.Subscriber{}, condition)
	if err != nil {
		return nil, err
	}

	return subscribers, nil
}

func (sb *subsRepository) Update(ctx context.Context, sub infoblog.Subscriber) error {
	return sb.crud.update(ctx, sub)
}

func NewSubscribeRepository(db *sqlx.DB) infoblog.SubscriberRepository {
	return &subsRepository{db: db, crud: crud{db: db}}
}
