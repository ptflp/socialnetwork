package request

//go:generate easytags $GOFILE

type PostCreateReq struct {
	Body string `json:"body"`
}
