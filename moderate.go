package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type Moderate struct {
	UUID      types.NullUUID  `json:"like_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	Type      types.NullInt64 `json:"type" db:"type" ops:"create" orm_type:"int" orm_default:"null"`
	Active    types.NullBool  `json:"active" db:"active" ops:"create" orm_type:"boolean" orm_default:"null"`
	CreatedAt time.Time       `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
}

func (l Moderate) OnCreate() string {
	return ""
}

func (l Moderate) TableName() string {
	return "moderates"
}

type ModerateRepository interface {
	Create(ctx context.Context, moderate Moderate) error
	Find(ctx context.Context, moderate Moderate) (Moderate, error)
	Update(ctx context.Context, moderate Moderate) error
	Delete(ctx context.Context, moderate Moderate) error
	List(ctx context.Context, limit, offset uint64) ([]Moderate, error)
	Listx(ctx context.Context, condition Condition) ([]Moderate, error)
}
