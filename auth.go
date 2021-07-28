package infoblog

import (
	"context"

	"gitlab.com/InfoBlogFriends/server/request"
)

type AuthService interface {
	SendCode(ctx context.Context, req *request.PhoneCodeRequest) bool
	CheckCode(ctx context.Context, req *request.CheckCodeRequest) (*request.AuthTokenData, error)
	EmailActivation(ctx context.Context, req *request.EmailActivationRequest) error
	EmailVerification(ctx context.Context, req *request.EmailVerificationRequest) (*request.AuthTokenData, error)
	SocialCallback(ctx context.Context, state string) (string, error)
	EmailLogin(ctx context.Context, req *request.EmailLoginRequest) (*request.AuthTokenData, error)
	RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*request.AuthTokenData, error)
}
