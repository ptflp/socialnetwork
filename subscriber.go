package infoblog

import (
	"context"
	"database/sql"
	"time"
)

type Subscriber struct {
	ID          int64
	UserID      int64
	SubscribeID int64
	Active      sql.NullBool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubscribesRepository interface {
	Create(ctx context.Context, uid, subID int64) (int64, error)
	FindByUser(ctx context.Context, uid int64) ([]Subscriber, error)
}
