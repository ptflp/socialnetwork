package infoblog

import "context"

type AuthService interface {
	SendCode(ctx context.Context, req *PhoneCodeRequest)
}
