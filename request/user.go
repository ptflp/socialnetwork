package request

//go:generate easytags $GOFILE

type UserSubscribeRequest struct {
	UUID string `json:"uuid"`
}
