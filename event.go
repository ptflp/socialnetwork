package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Event struct {
	UUID        types.NullUUID    `json:"event_id" db:"uuid" orm_type:"binary(16)" orm_default:"not null primary key"`
	Type        types.NullInt64   `json:"event_type" db:"type" orm_type:"int" orm_default:"default 1 not null"`
	ForeignUUID types.NullUUID    `json:"foreign_id" db:"foreign_uuid" ops:"update,create" orm_type:"binary(16)"`
	Shown       types.NullBool    `json:"shown" db:"shown" orm_type:"boolean" orm_default:"null"`
	Pushed      types.NullBool    `json:"pushed" db:"pushed" orm_type:"boolean" orm_default:"null"`
	UserUUID    types.NullUUID    `json:"user_id" db:"user_uuid" orm_type:"binary(16)" orm_default:"null" orm_index:"index"`
	ToUser      types.NullUUID    `json:"to_user" db:"to_user" orm_type:"binary(16)" orm_default:"null" orm_index:"index"`
	Price       types.NullFloat64 `json:"price" db:"price" orm_type:"decimal(13,4)" orm_default:"null" orm_index:"index"`
	Active      types.NullBool    `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
}

func (e Event) TableName() string {
	return "event"
}

func (e Event) OnCreate() string {
	return ""
}

type EventRepository interface {
	Create(ctx context.Context, event Event) error
	Find(ctx context.Context, event Event) (Event, error)
	Update(ctx context.Context, event Event) error
	List(ctx context.Context, limit, offset uint64) ([]Event, error)
}
