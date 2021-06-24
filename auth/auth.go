package auth

import (
	"context"
	"math/rand"
	"time"

	"gitlab.com/ptflp/infoblog-server/cache"

	"gitlab.com/ptflp/infoblog-server"
)

type AuthService struct {
	smsProvider    SmsProvider
	userRepository infoblog.UserRepository
	cache          cache.Cache
}

func NewAuthService(userRepository infoblog.UserRepository, cache cache.Cache) *AuthService {
	return &AuthService{userRepository: userRepository, cache: cache}
}

func (a *AuthService) SendCode(ctx context.Context, req *infoblog.PhoneCodeRequest) {
	code := genCode()
	err := a.smsProvider.SendCode(ctx, req.Phone, code)
	if err != nil {
		return
	}
}

func genCode() int {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(8999) + 1000

	return code
}

type SmsProvider interface {
	SendCode(ctx context.Context, phone string, code int) error
}

type ProviderMock struct {
	APIkey string
}

func (p *ProviderMock) SendCode(ctx context.Context, phone string, code int) error {
	return nil
}
