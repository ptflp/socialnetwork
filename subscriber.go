package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

//go:generate easytags $GOFILE

type Subscriber struct {
	ID             int64          `json:"id" db:"id"`
	UserUUID       types.NullUUID `json:"user_uuid" db:"user_uuid"`
	SubscriberUUID types.NullUUID `json:"subscriber_uuid" db:"subscriber_uuid"`
	Active         types.NullBool `json:"active" db:"active"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

func (s Subscriber) TableName() string {
	panic("implement me")
}

type SubscriberRepository interface {
	Create(ctx context.Context, sub Subscriber) (int64, error)
	FindByUser(ctx context.Context, user User) ([]Subscriber, error)
	Delete(ctx context.Context, sub Subscriber) error
	CountByUser(ctx context.Context, user User) (int64, error)
	CheckSubscribed(ctx context.Context, user User, subscriber User) bool
}
