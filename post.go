package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type PostEntity struct {
	UUID      types.NullUUID    `json:"post_id" db:"uuid" orm_type:"binary(16)" orm_default:"not null primary key" ops:"create"`
	Type      int64             `json:"post_type" db:"type" orm_type:"int" orm_default:"default 1 not null" ops:"create"`
	Body      string            `json:"description" db:"body" orm_type:"varchar(100)" orm_default:"not null" ops:"create,update"`
	UserUUID  types.NullUUID    `json:"user_id" db:"user_uuid" orm_type:"binary(16)" orm_default:"null" orm_index:"index" ops:"create"`
	Active    types.NullBool    `json:"active" db:"active" orm_type:"boolean" orm_default:"null" ops:"create"`
	Price     types.NullFloat64 `json:"price" db:"price" orm_type:"decimal(13,4)" orm_default:"null" orm_index:"index" ops:"create"`
	Likes     types.NullUint64  `json:"likes_count" db:"likes" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
	Views     types.NullUint64  `json:"views_count" db:"views" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
	CreatedAt time.Time         `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
	DeletedAt types.NullTime    `json:"deleted_at" db:"deleted_at" orm_type:"timestamp" orm_default:"null" orm_index:"index"`
}

func (p PostEntity) OnCreate() string {
	return ""
}

func (p PostEntity) TableName() string {
	return "posts"
}

type Post struct {
	PostEntity
	User          User   `json:"user" db:"user"`
	Files         []File `json:"files" db:"files"`
	CommentsCount int64  `json:"comments_count" db:"-"`
}

type PostRepository interface {
	Create(ctx context.Context, p Post) (int64, error)
	Update(ctx context.Context, p Post) error
	Delete(ctx context.Context, p Post) error
	Count(ctx context.Context, p Post, field, ops string) (Post, error)
	First(ctx context.Context) (Post, error)
	Listx(ctx context.Context, condition Condition) ([]PostEntity, error)

	Find(ctx context.Context, p Post) (Post, error)
	FindAll(ctx context.Context, user User, limit int64, offset int64) ([]Post, map[string]int, []string, error)
	FindAllRecent(ctx context.Context, limit, offset int64) ([]Post, map[string]int, []string, error)
	CountRecent(ctx context.Context) (int64, error)
	CountByUser(ctx context.Context, user User) (int64, error)
}
