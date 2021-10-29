package infoblog

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

//go:generate easytags $GOFILE

type Subscriber struct {
	UserUUID       types.NullUUID `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"null"`
	SubscriberUUID types.NullUUID `json:"subscriber_id" db:"subscriber_uuid" ops:"create" orm_type:"binary(16)" orm_default:"null"`
	Active         types.NullBool `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
	DeletedAt      sql.NullTime   `json:"deleted_at" db:"deleted_at" orm_type:"timestamp" orm_default:"null" orm_index:"index"`
}

func (s Subscriber) OnCreate() string {
	return "create unique index subscribes_user_uuid_subscriber_uuid_uindex on subscribes (user_uuid, subscriber_uuid);"
}

func (s Subscriber) TableName() string {
	return "subscribes"
}

type SubscriberRepository interface {
	Create(ctx context.Context, sub Subscriber) (int64, error)
	Update(ctx context.Context, sub Subscriber) error
	Updatex(ctx context.Context, sub Subscriber, condition Condition, ops string) error
	FindByUser(ctx context.Context, user User) ([]Subscriber, error)
	Delete(ctx context.Context, sub Subscriber) error
	CountByUser(ctx context.Context, user User) (int64, error)
	CheckSubscribed(ctx context.Context, user User, subscriber User) bool
	Listx(ctx context.Context, condition Condition) ([]Subscriber, error)
}
