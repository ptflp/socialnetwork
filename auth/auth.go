package auth

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"gitlab.com/InfoBlogFriends/server/email"
	"gitlab.com/InfoBlogFriends/server/hasher"

	"gitlab.com/InfoBlogFriends/server/validators"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/providers"

	"gitlab.com/InfoBlogFriends/server/session"

	"go.uber.org/zap"

	"gitlab.com/InfoBlogFriends/server/config"

	"gitlab.com/InfoBlogFriends/server/cache"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	EmailVerificationKey = "email:verification:%s"
)

type service struct {
	smsProvider    providers.SMS
	userRepository infoblog.UserRepository
	cache          cache.Cache
	config         *config.Config
	logger         *zap.Logger
	JWTKeys        *session.JWTKeys
	EmailClient    *email.Client
}

func NewAuthService(
	config *config.Config,
	userRepository infoblog.UserRepository,
	cache cache.Cache,
	logger *zap.Logger,
	keys *session.JWTKeys,
	smsProvider providers.SMS,
	emailClient *email.Client) *service {
	return &service{config: config, userRepository: userRepository, cache: cache, logger: logger, smsProvider: smsProvider, JWTKeys: keys, EmailClient: emailClient}
}

func (a *service) EmailActivation(ctx context.Context, req *request.EmailActivationRequest) error {
	// 1. Check user existance
	u, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err == nil && u.ID > 0 {
		return errors.New("user with specified email already exist")
	}

	activationUrl, activationID, err := a.generateActivationUrl(req.Email)
	if err != nil {
		return err
	}

	body, err := a.prepareEmailTemplate(activationUrl)
	if err != nil {
		return err
	}

	msg := email.NewMessage()
	msg.SetSubject("Активация учетной записи")
	msg.SetType(email.TypeHtml)
	msg.SetReceiver(req.Email)
	msg.SetBody(body)

	err = a.EmailClient.Send(msg)
	if err != nil {
		return err
	}

	data := req
	hashPass, err := hasher.HashPassword(req.Password)
	if err != nil {
		return err
	}
	data.Password = hashPass

	// 2. Set email code to cache
	a.cache.Set(fmt.Sprintf(EmailVerificationKey, activationID), data, 3*24*time.Hour)

	return nil
}

func (a *service) EmailVerification(ctx context.Context, req *request.EmailVerificationRequest) (string, error) {
	var u infoblog.User
	err := a.cache.Get(fmt.Sprintf(EmailVerificationKey, req.ActivationID), &u)
	if err != nil {
		return "", err
	}

	err = a.userRepository.CreateUserByEmailPassword(ctx, u.Email, u.Password)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			u, err = a.userRepository.FindByEmail(ctx, u.Email)
			if err != nil {
				return "", err
			}
			if u.EmailVerified == 1 {
				return "", fmt.Errorf("user with email %s already verified", u.Email)
			}
			if u.EmailVerified == 0 {
				u.EmailVerified = 1
				err = a.userRepository.Update(ctx, u)
				if err != nil {
					return "", err
				}
			}
		} else {
			return "", err
		}
	}

	u, err = a.userRepository.FindByEmail(ctx, u.Email)
	if err != nil {
		return "", err
	}
	if u.ID == 0 {
		return "", errors.New("email verification wrong user.ID")
	}

	token, err := a.JWTKeys.CreateToken(u)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *service) generateActivationUrl(email string) (string, string, error) {
	uid := uuid.NewV4()
	dh := uid.Bytes()
	dh = append(dh, []byte(email)...)
	hash := hex.EncodeToString(hasher.NewSHA256(dh))

	u, err := url.Parse(a.config.App.FrontEnd)
	if err != nil {
		return "", "", err
	}
	u.Path = fmt.Sprintf("email/%s", hash)

	return u.String(), hash, err
}

func (a *service) prepareEmailTemplate(activationUrl string) (bytes.Buffer, error) {
	tmpl, err := template.ParseFiles("./templates/email.html")
	if err != nil {
		a.logger.Error("email template parse", zap.Error(err))
		return bytes.Buffer{}, err
	}
	type EmailActivation struct {
		ActivationUrl string
	}

	b := bytes.Buffer{}

	err = tmpl.Execute(&b, EmailActivation{ActivationUrl: activationUrl})

	return b, err
}

func (a *service) SendCode(ctx context.Context, req *request.PhoneCodeRequest) bool {
	phone, err := validators.CheckPhoneFormat(req.Phone)
	if err != nil {
		return false
	}
	code := genCode()
	if a.config.SMSC.Dev {
		code = 3455
	}
	a.cache.Set("code:"+phone, &code, 15*time.Minute)
	if a.config.SMSC.Dev {
		return true
	}

	err = a.smsProvider.Send(ctx, phone, fmt.Sprintf("Ваш код: %d", code))
	if err != nil {
		a.logger.Error("send sms err", zap.String("phone", phone), zap.Int("code", code))
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
