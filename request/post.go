package request

//go:generate easytags $GOFILE

type PostCreateReq struct {
	Description string   `json:"description"`
	PostType    int64    `json:"post_type"`
	FilesID     []string `json:"files_id"`
	Price       *float64 `json:"price"`
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

type PostUUIDReq struct {
	UUID string `json:"post_id"`
}
