package infoblog

import (
	"context"
	"time"
)

type PostEntity struct {
	ID        int64     `json:"-" db:"id"`
	Type      int64     `json:"type" db:"type"`
	UUID      string    `json:"post_id" db:"uuid"`
	Body      string    `json:"body" db:"body"`
	FileID    int64     `json:"-" db:"file_id"`
	UserID    int64     `json:"-" db:"user_id"`
	UserUUID  string    `json:"user_id" db:"user_uuid"`
	Active    int64     `json:"active" db:"active"`
	FileUUID  string    `json:"file_id" db:"file_uuid"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
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
	FindAll(ctx context.Context, user User, limit int64, offset int64) ([]Post, map[int64]int, []int, error)
	FindAllRecent(ctx context.Context, limit, offset int64) ([]Post, map[int64]int, []int, error)
	CountRecent(ctx context.Context) (int64, error)
}
