package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Comments struct {
	UUID        types.NullUUID   `json:"uuid" db:"uuid" ops:"create" orm_type:"binary(16)"`
	ID          int64            `json:"id" db:"id" ops:"update,create" orm_type:"int"`
	Body        types.NullString `json:"body" db:"body" ops:"update,create" orm_type:"varchar(511)"`
	Type        int64            `json:"type" db:"type" ops:"update,create" orm_type:"int"`
	ForeignUUID types.NullUUID   `json:"foreign_id" db:"foreign_uuid" ops:"update,create" orm_type:"binary(16)"`
	UserUUID    types.NullUUID   `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
}

func (c Comments) TableName() string {
	return "comments"
}

type CommentsRepository interface {
	Update(ctx context.Context, comments Comments) error
	Find(ctx context.Context, comments Comments) (Comments, error)
	FindAll(ctx context.Context) ([]Comments, error)
	FindLimitOffset(ctx context.Context, limit, offset uint64) ([]Comments, error)
	CreateComments(ctx context.Context, comments Comments) error
}
