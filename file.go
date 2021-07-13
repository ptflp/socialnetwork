package infoblog

import (
	"context"
	"time"
)

type File struct {
	ID        int64     `json:"id" db:"id"`
	Type      int64     `json:"type" db:"type"`
	ForeignID int64     `json:"foreign_id" db:"foreign_id"`
	Dir       string    `json:"dir" db:"dir"`
	Name      string    `json:"name" db:"name"`
	Active    int64     `json:"active" db:"active"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type FileRepository interface {
	Create(ctx context.Context, p *File) (int64, error)
	Update(ctx context.Context, p File) error
	Delete(ctx context.Context, p File) error

	Find(ctx context.Context, id int64) (File, error)
	FindAll(ctx context.Context, postID int64) ([]File, error)
	FindByTypeFID(ctx context.Context, typeID int64, foreignID int64) ([]File, error)
}
