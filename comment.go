package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Comment struct {
	UUID        types.NullUUID   `json:"comment_id" db:"uuid" orm_type:"binary(16)" orm_default:"not null primary key" ops:"create"`
	Body        types.NullString `json:"body" db:"body" ops:"update,create" orm_type:"varchar(511)" orm_default:"null"`
	Type        int64            `json:"type" db:"type" ops:"update,create" orm_type:"int" orm_default:"not null"`
	ForeignUUID types.NullUUID   `json:"foreign_id" db:"foreign_uuid" ops:"update,create" orm_type:"binary(16)" orm_default:"null" orm_index:"index"`
	Active      types.NullBool   `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	UserUUID    types.NullUUID   `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null" orm_index:"index"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
	DeletedAt   types.NullTime   `json:"deleted_at" db:"deleted_at" orm_type:"timestamp" orm_default:"null" orm_index:"index" ops:"delete"`
}

func (c Comment) OnCreate() string {
	return ""
}

func (c Comment) TableName() string {
	return "comments"
}

type CommentsRepository interface {
	Create(ctx context.Context, comments Comment) error
	Find(ctx context.Context, comments Comment) (Comment, error)
	Update(ctx context.Context, comments Comment) error
	GetCount(ctx context.Context, condition Condition) (uint64, error)
	Delete(ctx context.Context, comments Comment) error
	List(ctx context.Context, limit, offset uint64) ([]Comment, error)
	Listx(ctx context.Context, condition Condition) ([]Comment, error)
}
