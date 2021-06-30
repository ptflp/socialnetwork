package request

type ProfileUpdateReq struct {
	Phone      *string `json:"phone"`
	Email      *string `json:"email"`
	Name       *string `json:"name"`
	SecondName *string `json:"second_name"`
}

type SetPasswordReq struct {
	Password string
}
