package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitlab.com/InfoBlogFriends/server/components"

	"github.com/google/uuid"

	"gitlab.com/InfoBlogFriends/server/email"
	"gitlab.com/InfoBlogFriends/server/hasher"

	"gitlab.com/InfoBlogFriends/server/validators"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/providers"

	"gitlab.com/InfoBlogFriends/server/session"

	"go.uber.org/zap"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	EmailVerificationKey = "email:verification:%s"
)

type service struct {
	smsProvider    providers.SMS
	userRepository infoblog.UserRepository
	components.Componenter
}

func NewAuthService(
	repositories infoblog.Repositories,
	cmps components.Componenter,
) *service {
	return &service{Componenter: cmps, userRepository: repositories.Users, smsProvider: cmps.SMS()}
}

func (a *service) EmailActivation(ctx context.Context, req *request.EmailActivationRequest) error {
	// 1. Check user existance
	u := infoblog.User{
		Email: infoblog.NewNullString(req.Email),
	}
	u, err := a.userRepository.FindByEmail(ctx, u)
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

	err = a.Email().Send(msg)
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
	a.Cache().Set(fmt.Sprintf(EmailVerificationKey, activationID), data, 3*24*time.Hour)

	return nil
}

func (a *service) EmailVerification(ctx context.Context, req *request.EmailVerificationRequest) (*request.AuthTokenData, error) {
	var u infoblog.User
	key := fmt.Sprintf(EmailVerificationKey, req.ActivationID)
	err := a.Cache().Get(key, &u)
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(89) + 10

	uUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	u.UUID = strings.Join([]string{uUUID.String(), fmt.Sprintf("-u%d", id)}, "")

	err = a.userRepository.CreateUserByEmailPassword(ctx, u)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			u, err = a.userRepository.FindByEmail(ctx, u)
			if err != nil {
				return nil, err
			}
			if u.EmailVerified.Bool == true {
				return nil, fmt.Errorf("user with email %s already verified", u.Email.String)
			}
			if u.EmailVerified.Bool == false {
				u.EmailVerified = infoblog.NewNullBool(true)
				err = a.userRepository.Update(ctx, u)
				if err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}

	u, err = a.userRepository.FindByEmail(ctx, u)
	if err != nil {
		return nil, err
	}
	if u.ID == 0 {
		return nil, errors.New("email verification wrong user.ID")
	}

	authTokens, err := a.JWTKeys().GenerateAuthTokens(&u)
	if err != nil {
		return nil, err
	}

	return authTokens, nil
}

func (a *service) RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*request.AuthTokenData, error) {
	var u infoblog.User
	refreshToken, err := a.JWTKeys().ExtractRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	key := strings.Join([]string{session.RefreshTokenKey, strconv.Itoa(int(refreshToken.UID)), refreshToken.Token}, ":")
	err = a.Cache().Get(key, &u)
	if err != nil {
		return nil, err
	}
	err = a.Cache().Del(key)
	if err != nil {
		a.Logger().Error("cache refresh_token del", zap.Error(err))
	}
	u, err = a.userRepository.Find(ctx, u)
	if err != nil {
		return nil, err
	}

	authTokens, err := a.JWTKeys().GenerateAuthTokens(&u)
	if err != nil {
		return nil, err
	}

	return authTokens, nil
}

func (a *service) EmailLogin(ctx context.Context, req *request.EmailLoginRequest) (*request.AuthTokenData, error) {
	var u infoblog.User
	u.Email = infoblog.NewNullString(req.Email)

	u, err := a.userRepository.FindByEmail(ctx, u)
	if err != nil {
		return nil, err
	}
	if u.ID == 0 {
		return nil, errors.New("wrong user.ID")
	}
	if !u.Password.Valid {
		return nil, errors.New("user password not set")
	}

	if !hasher.CheckPasswordHash(req.Password, u.Password.String) {
		return nil, errors.New("wrong email password")
	}

	token, err := a.JWTKeys().GenerateAuthTokens(&u)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *service) generateActivationUrl(email string) (string, string, error) {
	uid := uuid.New()
	dh, err := uid.MarshalBinary()
	if err != nil {
		return "", "", err
	}

	dh = append(dh, []byte(email)...)
	hash := hasher.NewSHA256(dh)

	u, err := url.Parse(a.Config().App.FrontEnd)
	if err != nil {
		return "", "", err
	}
	u.Path = fmt.Sprintf("email/%s", hash)

	return u.String(), hash, err
}

func (a *service) prepareEmailTemplate(activationUrl string) (bytes.Buffer, error) {
	tmpl, err := template.ParseFiles("./templates/email.html")
	if err != nil {
		a.Logger().Error("email template parse", zap.Error(err))
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
	if a.Config().SMSC.Dev {
		code = 3455
	}
	a.Cache().Set("code:"+phone, &code, 15*time.Minute)
	if a.Config().SMSC.Dev {
		return true
	}

	err = a.smsProvider.Send(ctx, phone, fmt.Sprintf("Ваш код: %d", code))
	if err != nil {
		a.Logger().Error("send sms err", zap.String("phone", phone), zap.Int("code", code))
	}

	return err == nil
}

func (a *service) CheckCode(ctx context.Context, req *request.CheckCodeRequest) (*request.AuthTokenData, error) {
	var code int
	phone, err := validators.CheckPhoneFormat(req.Phone)
	if err != nil {
		return nil, err
	}
	err = a.Cache().Get("code:"+phone, &code)
	if err != nil {
		return nil, err
	}

	if code != req.Code {
		return nil, errors.New("phone code mismatch")
	}

	phoneEnt := infoblog.NewNullString(phone)
	u := infoblog.User{
		Phone: phoneEnt,
	}
	u, err = a.userRepository.FindByPhone(ctx, u)
	u.Phone = phoneEnt
	if err != nil && err.Error() == "sql: no rows in result set" {
		rand.Seed(time.Now().UnixNano())
		id := rand.Intn(89) + 10

		uUUID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		u.UUID = strings.Join([]string{uUUID.String(), fmt.Sprintf("-u%d", id)}, "")

		err = a.userRepository.CreateUserByPhone(ctx, u)
		if err != nil {
			return nil, err
		}
		u, err = a.userRepository.Find(ctx, u)
		if err != nil {
			return nil, err
		}
	}

	token, err := a.JWTKeys().GenerateAuthTokens(&u)
	if err != nil {
		return nil, err
	}

	return token, err
}

func genCode() int {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(8999) + 1000

	return code
}
