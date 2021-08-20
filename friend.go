package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Friends struct {
	ID          types.NullInt64 `json:"id" db:"id" ops:"create" orm_type:"bigint" orm_default:"not null primary key"`
	FriendsUUID string          `json:"user_id" db:"user_uuid" ops:"create" orm_type:"varchar(40)" orm_default:"not null"`
	FriendUUID  string          `json:"friend_id" db:"friend_uuid" ops:"create" orm_type:"varchar(40)" orm_default:"not null"`
	Type        types.NullInt64 `json:"type" db:"type" ops:"update,create" orm_type:"int"`
	Active      types.NullBool  `json:"active" db:"active" ops:"create,update" orm_type:"boolean"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
}

func (f Friends) TableName() string {
	return "friends"
}

type FriendsRepository interface {
	Update(ctx context.Context, friends Friends) error
	Delete(ctx context.Context, friends Friends) error
	Find(ctx context.Context, friends Friends) (Friends, error)
	FindAll(ctx context.Context) ([]Friends, error)
	FindLimitOffset(ctx context.Context, limit, offset uint64) ([]Friends, error)
	CreateFriends(ctx context.Context, friends Friends) error
}
