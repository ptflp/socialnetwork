package infoblog

import (
	"gitlab.com/InfoBlogFriends/server/types"
	"time"
)

type Subscribe struct {
	UserUUID       types.NullUUID `json:"user_id" db:"user_uuid" ops:"create" orm_type:"binary(16)"`
	SubscriberUUID types.NullUUID `json:"subscriber_id" db:"subscriber_uuid" ops:"create" orm_type:"binary(16)"`
	Active         int64          `json:"active" db:"active" ops:"create,update"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

func (h Subscribe) TableName() string {
	return "subscribes"
}
