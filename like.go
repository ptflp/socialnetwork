package infoblog

import (
	"context"
	"time"
)

type Like struct {
	ID          int64     `json:"-" db:"id"`
	Type        int64     `json:"type" db:"type" ops:"create"`
	ForeignUUID string    `json:"foreign_id" db:"foreign_uuid" ops:"create"`
	UserUUID    string    `json:"user_id" db:"user_uuid" ops:"create"`
	LikerUUID   string    `json:"liker_id" db:"liker_uuid" ops:"create"`
	Active      NullBool  `json:"active" db:"active" ops:"create"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type LikeRepository interface {
	Upsert(ctx context.Context, like Like) error
	Find(ctx context.Context, like *Like) (Like, error)
}
