package request

import (
	"gitlab.com/InfoBlogFriends/server/types"
	"time"
)

//go:generate easytags $GOFILE

type ChatIDRequest struct {
	UUID *string `json:"chat_id"`
}

type ChatUpdateReq struct {
	UUID      types.NullUUID  `json:"chat_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	Type      types.NullInt64 `json:"type" db:"type" ops:"update,create" orm_type:"int"`
	CreatedAt time.Time       `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
}
