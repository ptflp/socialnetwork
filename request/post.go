package request

//go:generate easytags $GOFILE

type PostCreateReq struct {
	Description string   `json:"description"`
	PostType    int64    `json:"post_type"`
	FilesID     []string `json:"files_id"`
	Price       *float64 `json:"price"`
}

type PostUpdateReq struct {
	Body  string   `json:"description"`
	Price *float64 `json:"price"`
	UUID  string   `json:"post_id"`
}

type LimitOffsetReq struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type PostsFeedUserReq struct {
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
	UUID   string `json:"user_id"`
}

type PostUUIDReq struct {
	UUID string `json:"post_id"`
}

type UUIDReq struct {
	UUID string `json:"id"`
}

type LikeReq struct {
	UUID   string `json:"post_id"`
	Active bool   `json:"active"`
}

type CommentCreateReq struct {
	Body        string `json:"body"`
	ForeignUUID string `json:"post_id"`
}
