package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Files struct {
	UUID        types.NullUUID   `json:"uuid" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null"`
	Dir         types.NullString `json:"dir" db:"dir" ops:"update,create" orm_type:"varchar(100)" orm_default:"not null"`
	Active      types.NullInt64  `json:"active" db:"active" ops:"update,create" orm_type:"int" orm_default:"1"`
	Name        types.NullString `json:"name" db:"name" ops:"update,create" orm_type:"varchar(55)" orm_default:"not null"`
	Type        types.NullInt64  `json:"type" db:"type" ops:"update,create" orm_type:"int" orm_default:"not null"`
	ForeignUUID types.NullUUID   `json:"foreign_id" db:"foreign_uuid" ops:"update,create" orm_type:"binary(16)"`
	UserUUID    types.NullUUID   `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
}

func (f Files) TableName() string {
	return "files"
}

type FilesRepository interface {
	Update(ctx context.Context, files Files) error
	Find(ctx context.Context, files Files) (Files, error)
	FindAll(ctx context.Context) ([]Files, error)
	FindLimitOffset(ctx context.Context, limit, offset uint64) ([]Files, error)
	CreateFiles(ctx context.Context, files Files) error
}
