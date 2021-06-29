package auth

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"gitlab.com/InfoBlogFriends/server/session"

	"go.uber.org/zap"

	"gitlab.com/InfoBlogFriends/server/config"

	"gitlab.com/InfoBlogFriends/server/cache"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type service struct {
	smsProvider    SmsProvider
	userRepository infoblog.UserRepository
	cache          cache.Cache
	configApp      config.App
	logger         *zap.Logger
	JWTKeys        *session.JWTKeys
}

func NewAuthService(configApp config.App, userRepository infoblog.UserRepository, cache cache.Cache, logger *zap.Logger, keys *session.JWTKeys) *service {
	return &service{configApp: configApp, userRepository: userRepository, cache: cache, logger: logger, smsProvider: ProviderMock{}, JWTKeys: keys}
}

func (a *service) SendCode(ctx context.Context, req *infoblog.PhoneCodeRequest) bool {
	code := genCode()
	if a.configApp.Dev {
		code = 3455
	}
	a.cache.Set("code:"+req.Phone, &code, 15*time.Minute)
	err := a.smsProvider.SendCode(ctx, req.Phone, code)

	return err == nil
}

func (a *service) CheckCode(ctx context.Context, req *infoblog.CheckCodeRequest) (string, error) {
	var code int
	err := a.cache.Get("code:"+req.Phone, &code)
	if err != nil {
		return "", err
	}

	if code != req.Code {
		return "", errors.New("phone code mismatch")
	}

	u, err := a.userRepository.FindByPhone(ctx, req.Phone)
	if err != nil {
		err = a.userRepository.CreateUserByPhone(ctx, req.Phone)
		if err != nil {
			return "", err
		}
	}

	token, err := a.JWTKeys.CreateToken(u)
	if err != nil {
		return "", err
	}

	return token, err
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

func (p ProviderMock) SendCode(ctx context.Context, phone string, code int) error {
	return nil
}
