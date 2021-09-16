package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type File struct {
	UUID        types.NullUUID   `json:"file_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null"`
	Type        int64            `json:"type" db:"type" ops:"create,update" orm_type:"int" orm_default:"not null"`
	FileType    types.NullInt64  `json:"file_type" db:"file_type" ops:"create,update" orm_type:"int" orm_default:"null"`
	Status      types.NullInt64  `json:"status" db:"status" ops:"create,update" orm_type:"int" orm_default:"null"`
	Private     types.NullBool   `json:"private" db:"private" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	MimeType    types.NullString `json:"mime_type" db:"mime_type" ops:"create,update" orm_type:"varchar(34)" orm_default:"null"`
	ForeignUUID types.NullUUID   `json:"foreign_uuid" db:"foreign_uuid" ops:"create,update" orm_type:"binary(16)" orm_default:"null"`
	UserUUID    types.NullUUID   `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)" orm_default:"null"`
	Dir         string           `json:"dir" db:"dir" ops:"create,update" orm_type:"varchar(100)" orm_default:"not null"`
	Name        string           `json:"name" db:"name" ops:"create" orm_type:"varchar(50)" orm_default:"not null"`
	Active      int64            `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
}

func (f File) OnCreate() string {
	return ""
}

func (f File) TableName() string {
	return "files"
}

type FileRepository interface {
	Create(ctx context.Context, p *File) (int64, error)
	Update(ctx context.Context, p File) error
	UpdatePostUUID(ctx context.Context, ids []string, post Post) error
	UpdateFileType(ctx context.Context, file File, uuids ...types.NullUUID) error
	Delete(ctx context.Context, p File) error

	Find(ctx context.Context, f File) (File, error)
	FindAll(ctx context.Context, postUUID string) ([]File, error)
	FindByIDs(ctx context.Context, ids []string) ([]File, error)
	FindByTypeFUUID(ctx context.Context, typeID int64, foreignUUID string) ([]File, error)
	FindByPostsIDs(ctx context.Context, postsIDs []string) ([]File, error)
	Listx(ctx context.Context, condition Condition) ([]File, error)
}
