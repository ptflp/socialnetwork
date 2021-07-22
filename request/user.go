package request

//go:generate easytags $GOFILE

type UserIDRequest struct {
	UUID string `json:"user_id"`
}

type UserIDNickRequest struct {
	UUID     *string `json:"user_id"`
	NickName *string `json:"nickname"`
}

type PasswordRecoverRequest struct {
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

type CheckPhoneCodeRequest struct {
	Code  int64  `json:"code"`
	Phone string `json:"phone"`
}

type PasswordResetRequest struct {
	RecoverID string `json:"recover_id"`
	Password  string `json:"password"`
}
