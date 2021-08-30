package request

//go:generate easytags $GOFILE

type ModerateCreateReq struct {
	FilesID []string `json:"files_id"`
}

type ModerateUpdateStatusReq struct {
	Status int64
	Reason string
	UUID   string `json:"moderate_id" db:"uuid" ops:"create" orm_type:"binary(16)" orm_default:"not null primary key"`
}
