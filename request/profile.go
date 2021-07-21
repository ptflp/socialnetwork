package request

type ProfileUpdateReq struct {
	Phone          *string  `json:"phone,omitempty" db:"phone" ops:"update"`
	Email          *string  `json:"email,omitempty" db:"email" ops:"update"`
	Name           *string  `json:"name,omitempty" db:"name" ops:"update"`
	SecondName     *string  `json:"second_name,omitempty" db:"second_name" ops:"update"`
	Description    *string  `json:"description,omitempty" db:"description" ops:"update"`
	NickName       *string  `json:"nickname,omitempty" db:"nickname" ops:"update"`
	ShowSubs       *bool    `json:"show_subs,omitempty" db:"show_subs" ops:"update"`
	Cost           *float64 `json:"cost,omitempty" db:"cost" ops:"update"`
	Trial          *bool    `json:"trial,omitempty" db:"trial" ops:"update"`
	NotifyEmail    *bool    `json:"notify_email,omitempty" db:"notify_email" ops:"update"`
	NotifyTelegram *bool    `json:"notify_telegram,omitempty" db:"trial" ops:"update"`
	NotifyPush     *bool    `json:"notify_push,omitempty" db:"notify_push" ops:"update"`
	Language       *int64   `json:"language,omitempty" db:"language" ops:"update"`
}

type SetPasswordReq struct {
	Password    string  `json:"password"`
	OldPassword *string `json:"old_password"`
}
