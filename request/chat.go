package request

//go:generate easytags $GOFILE

type SendMessage struct {
	Message string `json:"message"`
	ToUUID  string `json:"to_uuid"`
}
