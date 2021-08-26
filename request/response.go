package request

import (
	"time"

	"gitlab.com/InfoBlogFriends/server/types"
)

//go:generate easytags $GOFILE
type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type AuthTokenResponse struct {
	Success bool          `json:"success"`
	Msg     string        `json:"msg"`
	Data    AuthTokenData `json:"data"`
}

type AuthTokenData struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserData `json:"user"`
}

type UserData struct {
	UUID           types.NullUUID    `json:"user_id" db:"uuid" ops:"create"`
	Phone          types.NullString  `json:"phone" db:"phone" ops:"update,create"`
	Email          types.NullString  `json:"email" db:"email" ops:"update,create"`
	Avatar         types.NullString  `json:"profile_image" db:"avatar" ops:"update"`
	Active         types.NullBool    `json:"active" db:"active" ops:"create"`
	Name           types.NullString  `json:"name" db:"name" ops:"update,create"`
	SecondName     types.NullString  `json:"second_name" db:"second_name" ops:"update,create"`
	EmailVerified  types.NullBool    `json:"email_verified" db:"email_verified"`
	Description    types.NullString  `json:"description" db:"description" ops:"update,create"`
	NickName       types.NullString  `json:"nickname" db:"nickname" ops:"update,create"`
	ShowSubs       types.NullBool    `json:"show_subs" db:"show_subs" ops:"update,create"`
	Cost           types.NullFloat64 `json:"cost" db:"cost" ops:"update,create"`
	Trial          types.NullBool    `json:"trial" db:"trial" ops:"update,create"`
	NotifyEmail    types.NullBool    `json:"notify_email" db:"notify_email" ops:"update,create"`
	NotifyTelegram types.NullBool    `json:"notify_telegram" db:"notify_telegram" ops:"update,create"`
	NotifyPush     types.NullBool    `json:"notify_push" db:"notify_push" ops:"update,create"`
	Language       int64             `json:"language" db:"language" ops:"update,create"`
	Likes          types.NullUint64  `json:"likes_count" db:"likes" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
	Subscribes     types.NullUint64  `json:"subscribes_count" db:"subscribes" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
	Subscribers    types.NullUint64  `json:"subscribers_count" db:"subscribers" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`

	IsFriends    bool `json:"is_friends"`
	IsSubscriber bool `json:"is_subscriber"`

	PasswordSet *bool `json:"password_set,omitempty"`
	AvatarSet   bool  `json:"profile_image_set"`

	Counts *UserDataCounts `json:"counts,omitempty"`
}

type UserDataCounts struct {
	Posts       int64 `json:"posts"`
	Subscribers int64 `json:"subscribers"`
	Friends     int64 `json:"friends"`
	Likes       int64 `json:"likes"`
}

type PostDataResponse struct {
	UUID   string           `json:"post_id"`
	Body   string           `json:"description"`
	Type   int64            `json:"post_type"`
	Files  []PostFileData   `json:"files"`
	User   UserData         `json:"user"`
	Price  float64          `json:"price"`
	Counts PostCountData    `json:"counts"`
	Likes  types.NullUint64 `json:"likes_count" db:"likes" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
	Views  types.NullUint64 `json:"views_count" db:"views" orm_type:"bigint unsigned" orm_default:"null" orm_index:"index" ops:"count"`
}

type PostFileData struct {
	Link string `json:"link"`
	UUID string `json:"file_id"`
}

type PostCountData struct {
	Likes    int64 `json:"likes"`
	Comments int64 `json:"comments"`
}

type PostsFeedData struct {
	Count uint64             `json:"count"`
	Posts []PostDataResponse `json:"posts"`
}

type PostsFeedResponse struct {
	Success bool          `json:"success"`
	Msg     string        `json:"msg"`
	Data    PostsFeedData `json:"data"`
}

type RecoverChekPhoneResponse struct {
	Success bool                  `json:"success"`
	Msg     string                `json:"msg"`
	Data    RecoverCheckPhoneData `json:"data"`
}

type RecoverCheckPhoneData struct {
	RecoverID string `json:"recover_id"`
}

type ChatData struct {
	UUID      types.NullUUID  `json:"chat_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
	Type      types.NullInt64 `json:"type" db:"type" ops:"update,create" orm_type:"int"`
	CreatedAt time.Time       `json:"created_at" db:"created_at" orm_type:"timestamp" orm_default:"default (now()) not null" orm_index:"index"`
}
