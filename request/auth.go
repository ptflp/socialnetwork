package request

type PhoneCodeRequest struct {
	Phone string `json:"phone"`
}

type CheckCodeRequest struct {
	Phone string `json:"phone"`
	Code  int
}

type EmailActivationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailVerificationRequest struct {
	ActivationID string `json:"activation_id"`
}
