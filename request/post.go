package request

//go:generate easytags $GOFILE

type PostCreateReq struct {
	Body string `json:"body"`
}

type PostsFeedReq struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}
