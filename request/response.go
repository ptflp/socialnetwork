package request

import "time"

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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserData struct {
	UUID           string    `json:"user_id" db:"uuid" ops:"create"`
	Phone          string    `json:"phone" db:"phone" ops:"update,create"`
	Email          string    `json:"email" db:"email" ops:"update,create"`
	Password       string    `json:"password,omitempty" db:"password" ops:"create"`
	Active         bool      `json:"active" db:"active" ops:"create"`
	Name           string    `json:"name" db:"name" ops:"update,create"`
	SecondName     string    `json:"second_name" db:"second_name" ops:"update,create"`
	EmailVerified  bool      `json:"email_verified" db:"email_verified"`
	Description    string    `json:"description" db:"description" ops:"update,create"`
	NickName       string    `json:"nickname" db:"nickname" ops:"update,create"`
	ShowSubs       bool      `json:"show_subs" db:"show_subs" ops:"update,create"`
	Cost           float64   `json:"cost" db:"cost" ops:"update,create"`
	Trial          bool      `json:"trial" db:"trial" ops:"update,create"`
	NotifyEmail    bool      `json:"notify_email" db:"notify_email" ops:"update,create"`
	NotifyTelegram bool      `json:"notify_telegram" db:"notify_telegram" ops:"update,create"`
	NotifyPush     bool      `json:"notify_push" db:"notify_push" ops:"update,create"`
	Language       int64     `json:"language" db:"language" ops:"update,create"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type PostDataResponse struct {
	UUID   string         `json:"post_id"`
	Body   string         `json:"body"`
	Files  []PostFileData `json:"files"`
	User   UserData       `json:"user"`
	Counts PostCountData  `json:"counts"`
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
	Count int64              `json:"count"`
	Posts []PostDataResponse `json:"posts"`
}

type PostsFeedResponse struct {
	Success bool          `json:"success"`
	Msg     string        `json:"msg"`
	Data    PostsFeedData `json:"data"`
}
