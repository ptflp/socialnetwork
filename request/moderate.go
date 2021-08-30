package request

//go:generate easytags $GOFILE

type ModerateCreateReq struct {
	FilesID []string `json:"files_id"`
}
