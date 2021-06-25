package infoblog

import (
	"context"
)

type AuthService interface {
	SendCode(ctx context.Context, req *PhoneCodeRequest) bool
	CheckCode(ctx context.Context, req *CheckCodeRequest) (string, error)
}
