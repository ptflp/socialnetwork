package infoblog

type PhoneCodeRequest struct {
	Phone string
}

type CheckCodeRequest struct {
	Phone string
	Code  int
}
