package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"gitlab.com/InfoBlogFriends/server/config"

	"gitlab.com/InfoBlogFriends/server/cache"

	jwt "github.com/dgrijalva/jwt-go"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	privKeyPath = "./keys/private"
	pubKeyPath  = "./keys/public"
)

var signBytes []byte
var verifyBytes []byte

type AuthService struct {
	smsProvider    SmsProvider
	userRepository infoblog.UserRepository
	cache          cache.Cache
	configApp      config.App
	logger         *zap.Logger
}

func NewAuthService(userRepository infoblog.UserRepository, cache cache.Cache, logger *zap.Logger) *AuthService {
	return &AuthService{userRepository: userRepository, cache: cache, logger: logger, smsProvider: ProviderMock{}}
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

func (a *AuthService) CheckCode(ctx context.Context, req *infoblog.CheckCodeRequest) (string, error) {
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

	token, err := a.CreateToken(u)
	if err != nil {
		return "", err
	}

	return token, err
}

func (a *AuthService) CreateToken(u infoblog.User) (string, error) {
	signKey, _, _ := a.ReadKeys()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"ExpiresAt": time.Now().UTC().Add(time.Minute * 20).Unix(),
		"ID":        u.ID,
		"Email":     u.Email,
		"Phone":     u.Phone,
	})

	return token.SignedString(signKey)
}

func (a *AuthService) ReadKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	var err error
	var signKey *rsa.PrivateKey
	var verifyKey *rsa.PublicKey

	if signBytes == nil {
		signBytes, err = ioutil.ReadFile(privKeyPath)
		if err != nil {
			a.logger.Error("read private key err", zap.Error(err))
		}
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		a.logger.Error("read private key err", zap.Error(err))
	}

	if verifyBytes == nil {
		verifyBytes, err = ioutil.ReadFile(pubKeyPath)
		if err != nil {
			a.logger.Error("read private key err", zap.Error(err))
		}
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		a.logger.Error("read private key err", zap.Error(err))
	}
	return signKey, verifyKey, err
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
