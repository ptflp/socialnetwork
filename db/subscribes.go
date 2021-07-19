package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	createSubscribe = "INSERT INTO subscribes (user_uuid, subscriber_uuid, active) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE active = 1"
	deleteSubscribe = "INSERT INTO subscribes (user_uuid, subscriber_uuid, active) VALUES (?, ?, 0) ON DUPLICATE KEY UPDATE active = 0"
)

type subsRepository struct {
	db *sqlx.DB
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

func NewSubscribeRepository(db *sqlx.DB) infoblog.SubscriberRepository {
	return &subsRepository{db: db}
}
