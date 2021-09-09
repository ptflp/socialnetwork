package request

//go:generate easytags $GOFILE

type SendMessageReq struct {
	Message  string `json:"message"`
	ChatUUID string `json:"chat_id"`
}

type GetInfoReq struct {
	UserUUID *string `json:"user_id"`
	ChatUUID *string `json:"chat_id"`
}
