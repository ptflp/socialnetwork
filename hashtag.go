package infoblog

import (
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type HashTag struct {
	ID        int64          `json:"-" db:"id" ops:"create" orm_type:"bigint" orm_default:"not null primary key"`
	Tag       string         `json:"tag" db:"tag" ops:"create,update" orm_type:"varchar(255)" orm_default:"not null"`
	Active    types.NullBool `json:"active" db:"active" ops:"create,update" orm_type:"boolean" orm_default:"null"`
	CreatedAt time.Time      `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
}

func (h HashTag) OnCreate() string {
	return ""
}

func (h HashTag) TableName() string {
	return "hashtags"
}
