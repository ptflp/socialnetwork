package infoblog

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type PostEntity struct {
	UUID      types.NullUUID    `json:"post_id" db:"uuid" gorm:"primaryKey,type:binary(16)"`
	Type      int64             `json:"post_type" db:"type" gorm:"default:1"`
	Body      string            `json:"description" db:"body" gorm:"type:varchar(100)"`
	UserUUID  types.NullUUID    `json:"user_id" db:"user_uuid" gorm:"type:binary(16),index"`
	Active    types.NullBool    `json:"active" db:"active"`
	Price     types.NullFloat64 `json:"price" db:"price" gorm:"type:decimal(13,4)"`
	Likes     types.NullInt64   `json:"likes" db:"likes" gorm:"index"`
	CreatedAt time.Time         `json:"created_at" db:"created_at" gorm:"index,type:timestamp"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at" gorm:"index,type:timestamp"`
	DeletedAt sql.NullTime      `json:"deleted_at" db:"deleted_at" gorm:"index"`
}

func (p *PostEntity) TableName() string {
	return "posts"
}

type Post struct {
	PostEntity
	User          User   `json:"user" db:"user"`
	Files         []File `json:"files" db:"files"`
	Likes         []Like `json:"likes" db:"likes"`
	LikesCount    int64  `json:"likes_count" db:"-"`
	CommentsCount int64  `json:"comments_count" db:"-"`
}

type PostRepository interface {
	Create(ctx context.Context, p Post) (int64, error)
	Update(ctx context.Context, p Post) error
	Delete(ctx context.Context, p Post) error

	Find(ctx context.Context, p Post) (Post, error)
	FindAll(ctx context.Context, user User, limit int64, offset int64) ([]Post, map[string]int, []string, error)
	FindAllRecent(ctx context.Context, limit, offset int64) ([]Post, map[string]int, []string, error)
	CountRecent(ctx context.Context) (int64, error)
	CountByUser(ctx context.Context, user User) (int64, error)
}
