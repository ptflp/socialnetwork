package request

type ProfileUpdateReq struct {
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	NickName    *string `json:"nickname"`
}

type SetPasswordReq struct {
	Password string
}
