package infoblog

import (
	"context"

	"gitlab.com/InfoBlogFriends/server/handlers"
)

type AuthService interface {
	SendCode(ctx context.Context, req *handlers.PhoneCodeRequest) bool
	CheckCode(ctx context.Context, req *handlers.CheckCodeRequest) (string, error)
}
