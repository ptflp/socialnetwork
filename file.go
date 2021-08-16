package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type File struct {
	UUID        types.NullUUID `json:"file_id" db:"uuid" ops:"create"`
	Type        int64          `json:"type" db:"type" ops:"create,update"`
	ForeignUUID types.NullUUID `json:"foreign_uuid" db:"foreign_uuid" ops:"create,update"`
	UserUUID    types.NullUUID `json:"user_id" db:"user_uuid" ops:"create"`
	Dir         string         `json:"dir" db:"dir" ops:"create,update"`
	Name        string         `json:"name" db:"name" ops:"create"`
	Active      int64          `json:"active" db:"active" ops:"create,update"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

type FileRepository interface {
	Create(ctx context.Context, p *File) (int64, error)
	Update(ctx context.Context, p File) error
	UpdatePostUUID(ctx context.Context, ids []string, post Post) error
	Delete(ctx context.Context, p File) error

	Find(ctx context.Context, f File) (File, error)
	FindAll(ctx context.Context, postUUID string) ([]File, error)
	FindByIDs(ctx context.Context, ids []string) ([]File, error)
	FindByTypeFUUID(ctx context.Context, typeID int64, foreignUUID string) ([]File, error)
	FindByPostsIDs(ctx context.Context, postsIDs []string) ([]File, error)
}
