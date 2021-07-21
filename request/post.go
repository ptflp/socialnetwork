package request

//go:generate easytags $GOFILE

type PostCreateReq struct {
	Body    string   `json:"body"`
	FilesID []string `json:"files_id"`
}

type PostsFeedReq struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type PostsFeedUserReq struct {
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
	UUID   string `json:"user_id"`
}

type PostLikeReq struct {
	UUID string `json:"post_id"`
}
