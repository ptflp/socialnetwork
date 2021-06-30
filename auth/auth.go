package auth

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/InfoBlogFriends/server/validators"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/providers"

	"gitlab.com/InfoBlogFriends/server/session"

	"go.uber.org/zap"

	"gitlab.com/InfoBlogFriends/server/config"

	"gitlab.com/InfoBlogFriends/server/cache"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type service struct {
	smsProvider    providers.SMS
	userRepository infoblog.UserRepository
	cache          cache.Cache
	config         *config.Config
	logger         *zap.Logger
	JWTKeys        *session.JWTKeys
}

func NewAuthService(
	config *config.Config,
	userRepository infoblog.UserRepository,
	cache cache.Cache,
	logger *zap.Logger,
	keys *session.JWTKeys,
	smsProvider providers.SMS) *service {
	return &service{config: config, userRepository: userRepository, cache: cache, logger: logger, smsProvider: smsProvider, JWTKeys: keys}
}

func (a *service) SendCode(ctx context.Context, req *request.PhoneCodeRequest) bool {
	code := genCode()
	if a.config.SMSC.Dev {
		code = 3455
	}
	a.cache.Set("code:"+req.Phone, &code, 15*time.Minute)
	if a.config.SMSC.Dev {
		return true
	}

	err := a.smsProvider.Send(ctx, req.Phone, fmt.Sprintf("Ваш код: %d", code))
	if err != nil {
		a.logger.Error("send sms err", zap.String("phone", req.Phone), zap.Int("code", code))
	}

	return err == nil
}

func (a *service) CheckCode(ctx context.Context, req *request.CheckCodeRequest) (string, error) {
	var code int
	phone, err := validators.CheckPhoneFormat(req.Phone)
	if err != nil {
		return "", err
	}
	err = a.cache.Get("code:"+phone, &code)
	if err != nil {
		return "", err
	}

	if code != req.Code {
		return "", errors.New("phone code mismatch")
	}

	u, err := a.userRepository.FindByPhone(ctx, phone)
	if err != nil {
		err = a.userRepository.CreateUserByPhone(ctx, phone)
		if err != nil {
			return "", err
		}
		u, err = a.userRepository.FindByPhone(ctx, phone)
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
