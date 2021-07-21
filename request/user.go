package request

//go:generate easytags $GOFILE

type UserIDRequest struct {
	UUID string `json:"user_id"`
}

type UserIDNickRequest struct {
	UUID     *string `json:"user_id"`
	NickName *string `json:"nickname"`
}
