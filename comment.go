package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Comments struct {
	UUID        types.NullUUID   `json:"comment_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null"`
	ID          int64            `json:"id" db:"id" ops:"update,create" orm_type:"int" orm_default:"not null"`
	Body        types.NullString `json:"body" db:"body" ops:"update,create" orm_type:"varchar(511)" orm_default:"null"`
	Type        int64            `json:"type" db:"type" ops:"update,create" orm_type:"int" orm_default:"not null"`
	ForeignUUID types.NullUUID   `json:"foreign_id" db:"foreign_uuid" ops:"update,create" orm_type:"binary(16)" orm_default:"null"`
	Active      types.NullBool   `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	UserUUID    types.NullUUID   `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
}

func (c Comments) TableName() string {
	return "comments"
}

type CommentsRepository interface {
	Create(ctx context.Context, comments Comments) error
	Find(ctx context.Context, comments Comments) (Comments, error)
	Update(ctx context.Context, comments Comments) error
	Delete(ctx context.Context, comments Comments) error
	List(ctx context.Context, limit, offset uint64) ([]Comments, error)
}
