package infoblog

import (
	"context"
	"time"
)

type Post struct {
	ID        int64     `json:"id" db:"id"`
	Body      string    `json:"body" db:"body"`
	FileID    int64     `json:"file_id" db:"file_id"`
	UserID    int64     `json:"user_id"`
	Active    int64     `json:"active" db:"active"`
	Type      int64     `json:"type" db:"type""`
	UUID      string    `json:"uuid" db:"uuid"`
	FileUUID  string    `json:"file_uuid" db:"file_uuid"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PostRepository interface {
	Create(ctx context.Context, p Post) (int64, error)
	Update(ctx context.Context, p Post) error
	Delete(ctx context.Context, p Post) error

	Find(ctx context.Context, id int64) (Post, error)
	FindAll(ctx context.Context, uid int64) ([]Post, error)
	FindAllRecent(ctx context.Context, limit, offset int64) ([]Post, error)
}
