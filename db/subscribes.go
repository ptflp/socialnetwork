package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	createSubscribe = "INSERT INTO subscribes (user_id, subscribe_id, active) VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE active = 1"
)

type subsRepository struct {
	db *sqlx.DB
}

func (sub *subsRepository) Create(ctx context.Context, uid, subID int64) (int64, error) {
	res, err := sub.db.ExecContext(ctx, createSubscribe, uid, subID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (sub *subsRepository) FindByUser(ctx context.Context, uid int64) ([]infoblog.Subscriber, error) {
	panic("implement me")
}

func NewSubscribeRepository(db *sqlx.DB) infoblog.SubscribesRepository {
	return &subsRepository{db: db}
}
