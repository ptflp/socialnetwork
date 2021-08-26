package infoblog

import (
	"context"
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

type User struct {
	UUID           types.NullUUID    `json:"user_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	Role           types.NullInt64   `json:"role" db:"role" orm_type:"int" orm_default:"null" ops:"update"`
	Phone          types.NullString  `json:"phone" db:"phone" ops:"update,create" orm_type:"varchar(34)" orm_index:"index,unique"`
	Email          types.NullString  `json:"email" db:"email" ops:"update,create" orm_type:"varchar(89)" orm_index:"index,unique"`
	Avatar         types.NullString  `json:"profile_image" db:"avatar" ops:"update" orm_type:"varchar(144)"`
	Password       types.NullString  `json:"password,omitempty" db:"password" ops:"create" orm_type:"varchar(60)"`
	Active         types.NullBool    `json:"active" db:"active" ops:"create,update" orm_type:"boolean"`
	Name           types.NullString  `json:"name" db:"name" ops:"update,create" orm_type:"varchar(55)"`
	SecondName     types.NullString  `json:"second_name" db:"second_name" ops:"update,create" orm_type:"varchar(55)"`
	EmailVerified  types.NullBool    `json:"email_verified" db:"email_verified" orm_type:"boolean"`
	Description    types.NullString  `json:"description" db:"description" ops:"update,create" orm_type:"varchar(233)"`
	NickName       types.NullString  `json:"nickname" db:"nickname" ops:"update,create" orm_type:"varchar(30)"`
	ShowSubs       types.NullBool    `json:"show_subs" db:"show_subs" ops:"update,create" orm_type:"boolean"`
	Cost           types.NullFloat64 `json:"cost" db:"cost" ops:"update,create" orm_type:"decimal(13,4)"`
	Trial          types.NullBool    `json:"trial" db:"trial" ops:"update,create" orm_type:"boolean"`
	NotifyEmail    types.NullBool    `json:"notify_email" db:"notify_email" ops:"update,create" orm_type:"boolean"`
	NotifyTelegram types.NullBool    `json:"notify_telegram" db:"notify_telegram" ops:"update,create" orm_type:"boolean"`
	NotifyPush     types.NullBool    `json:"notify_push" db:"notify_push" ops:"update,create" orm_type:"boolean"`
	Language       types.NullInt64   `json:"language" db:"language" ops:"update,create" orm_type:"int"`
	FacebookID     types.NullInt64   `json:"facebook_id" db:"facebook_id" ops:"update,create" orm_type:"bigint unsigned"`
	GoogleID       types.NullString  `json:"google_id" db:"google_id" ops:"update,create" orm_type:"varchar(21)"`
	Likes          types.NullUint64  `json:"likes_count" db:"likes" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
	Subscribes     types.NullUint64  `json:"subscribes_count" db:"subscribes" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
	Subscribers    types.NullUint64  `json:"subscribers_count" db:"subscribers" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at" orm_type:"timestamp" orm_default:"default (now()) null on update CURRENT_TIMESTAMP" orm_index:"index"`
	DeletedAt      types.NullTime    `json:"deleted_at" db:"deleted_at" orm_type:"timestamp" orm_default:"null" orm_index:"index"`
}

func (u User) OnCreate() string {
	return ""
}

func (u User) TableName() string {
	return "users"
}

type UserRepository interface {
	Update(ctx context.Context, user User) error
	SetPassword(ctx context.Context, user User) error

	Find(ctx context.Context, user User) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindLimitOffset(ctx context.Context, limit, offset uint64) ([]User, error)
	FindByPhone(ctx context.Context, user User) (User, error)
	FindByEmail(ctx context.Context, user User) (User, error)
	FindByNickname(ctx context.Context, user User) (User, error)
	FindLikeNickname(ctx context.Context, nickname string) ([]User, error)
	FindByFacebook(ctx context.Context, user User) (User, error)
	FindByGoogle(ctx context.Context, user User) (User, error)
	Count(ctx context.Context, user User, field, ops string) (User, error)

	CreateUser(ctx context.Context, user User) error
	CreateUserByEmailPassword(ctx context.Context, user User) error

	Listx(ctx context.Context, condition Condition) ([]User, error)
}
