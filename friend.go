package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Friend struct {
	ID          types.NullInt64 `json:"id" db:"id" ops:"create" orm_type:"bigint" orm_default:"not null primary key"`
	FriendsUUID string          `json:"user_id" db:"user_uuid" ops:"create" orm_type:"varchar(40)" orm_default:"not null"`
	FriendUUID  string          `json:"friend_id" db:"friend_uuid" ops:"create" orm_type:"varchar(40)" orm_default:"not null"`
	Type        types.NullInt64 `json:"type" db:"type" ops:"update,create" orm_type:"int" orm_default:"null"`
	Active      types.NullBool  `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
}

func (f Friend) OnCreate() string {
	return ""
}

func (f Friend) TableName() string {
	return "friends"
}

type FriendRepository interface {
	Create(ctx context.Context, friend Friend) error
	Find(ctx context.Context, friend Friend) (Friend, error)
	Update(ctx context.Context, friend Friend) error
	Delete(ctx context.Context, friend Friend) error
	List(ctx context.Context, limit, offset uint64) ([]Friend, error)
}
