package auth

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"gitlab.com/InfoBlogFriends/server/config"

	"gitlab.com/InfoBlogFriends/server/cache"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type AuthService struct {
	smsProvider    SmsProvider
	userRepository infoblog.UserRepository
	cache          cache.Cache
	configApp      config.App
	logger         *zap.Logger
}

func NewAuthService(userRepository infoblog.UserRepository, cache cache.Cache, logger *zap.Logger) *AuthService {
	return &AuthService{userRepository: userRepository, cache: cache, logger: logger}
}

func (a *AuthService) SendCode(ctx context.Context, req *infoblog.PhoneCodeRequest) bool {
	code := genCode()
	a.cache.Set("code:"+req.Phone, &code, 15*time.Minute)
	err := a.smsProvider.SendCode(ctx, req.Phone, code)
	if err != nil {
		return false
	}

	return true
}

func (a *AuthService) CheckCode(ctx context.Context, req *infoblog.CheckCodeRequest) bool {
	var code int
	err := a.cache.Get("code:"+req.Phone, &code)
	if err != nil {
		return false
	}

	_, err = a.userRepository.FindByPhone(ctx, req.Phone)
	if err != nil {
		err = a.userRepository.CreateUserByPhone(ctx, req.Phone)
		if err != nil {

		}
	}

	return true
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
