package infoblog

import (
	"context"

	"gitlab.com/InfoBlogFriends/server/request"
)

type AuthService interface {
	SendCode(ctx context.Context, req *request.PhoneCodeRequest) bool
	CheckCode(ctx context.Context, req *request.CheckCodeRequest) (string, error)
	EmailActivation(ctx context.Context, req *request.EmailActivationRequest) error
	EmailVerification(ctx context.Context, req *request.EmailVerificationRequest) (string, error)
}
