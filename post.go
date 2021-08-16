package infoblog

import (
	"context"
	"time"
)

type PostEntity struct {
	ID        int64       `json:"-" db:"id"`
	Type      int64       `json:"post_type" db:"type"`
	UUID      NullUUID    `json:"post_id" db:"uuid"`
	Body      string      `json:"description" db:"body"`
	UserUUID  NullUUID    `json:"user_id" db:"user_uuid"`
	Active    int64       `json:"active" db:"active"`
	Price     NullFloat64 `json:"price" db:"price"`
	FileUUID  NullUUID    `json:"file_id" db:"file_uuid"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
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
