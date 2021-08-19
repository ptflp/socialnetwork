package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Like struct {
	Type        int64          `json:"type" db:"type" ops:"create"`
	ForeignUUID types.NullUUID `json:"foreign_id" db:"foreign_uuid" ops:"create"`
	UserUUID    types.NullUUID `json:"user_id" db:"user_uuid" ops:"create"`
	LikerUUID   types.NullUUID `json:"liker_id" db:"liker_uuid" ops:"create"`
	Active      types.NullBool `json:"active" db:"active" ops:"create"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

func (l Like) TableName() string {
	return "likes"
}

type LikeRepository interface {
	Upsert(ctx context.Context, like Like) error
	Find(ctx context.Context, like *Like) (Like, error)
	CountByUser(ctx context.Context, user User) (int64, error)
	CountByPost(ctx context.Context, like Like) (int64, error)
}
