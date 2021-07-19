package request

//go:generate easytags $GOFILE

type UserSubscriberRequest struct {
	UUID string `json:"user_id"`
}
