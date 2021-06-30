package request

type PhoneCodeRequest struct {
	Phone string `json:"phone"`
}

type CheckCodeRequest struct {
	Phone string
	Code  int
}
