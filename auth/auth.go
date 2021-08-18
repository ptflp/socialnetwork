package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"gitlab.com/InfoBlogFriends/server/types"

	"gitlab.com/InfoBlogFriends/server/decoder"

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
	PhoneRegistrationKey = "phone:registration:%s"

	SocialsAuthKey = "socials:auth:%s"
)

type Provider struct{}
type State struct{}

type service struct {
	*decoder.Decoder
	smsProvider    providers.SMS
	userRepository infoblog.UserRepository
	components.Componenter
}

func emailBody(emailTo, subject, data string) []byte {
	address := "friends@22byte.com"
	name := "friends@22byte.com"
	from := mail.NewEmail(name, address)
	address = emailTo
	name = emailTo
	to := mail.NewEmail(name, address)
	content := mail.NewContent(email.TypeHtml, data)
	m := mail.NewV3MailInit(from, subject, to, content)

	return mail.GetRequestBody(m)
}

func NewAuthService(
	repositories infoblog.Repositories,
	cmps components.Componenter,
) *service {
	return &service{Componenter: cmps, userRepository: repositories.Users, smsProvider: cmps.SMS(), Decoder: decoder.NewDecoder()}
}

func (a *service) EmailActivation(ctx context.Context, req *request.EmailActivationRequest) error {
	// 1. Check user existance
	u := infoblog.User{
		Email: types.NewNullString(req.Email),
	}
	u, err := a.userRepository.FindByEmail(ctx, u)
	if err == nil && !u.UUID.Valid {
		return errors.New("user with specified email already exist")
	}

	activationUrl, activationID, err := a.generateActivationUrl(req.Email)
	if err != nil {
		return err
	}

	html, err := a.prepareEmailTemplate(activationUrl)
	if err != nil {
		return err
	}

	sendEmailReq := sendgrid.GetRequest("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", "/v3/mail/send", "https://api.sendgrid.com")
	sendEmailReq.Method = "POST"
	body := emailBody(req.Email, "Активация учетной записи", html.String())
	sendEmailReq.Body = body
	response, err := sendgrid.API(sendEmailReq)
	if err != nil {
		return err
	}
	_ = response

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
	u, err := createDefaultUser()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf(EmailVerificationKey, req.ActivationID)
	err = a.Cache().Get(key, &u)
	if err != nil {
		return nil, err
	}

	u.EmailVerified = types.NewNullBool(true)
	err = a.userRepository.CreateUser(ctx, u)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			u, err = a.userRepository.FindByEmail(ctx, u)
			if err != nil {
				return nil, err
			}
			if u.EmailVerified.Bool {
				return nil, fmt.Errorf("user with email %s already verified", u.Email.String)
			}
		} else {
			return nil, err
		}
	}

	u, err = a.userRepository.FindByEmail(ctx, u)
	if err != nil {
		return nil, err
	}
	if !u.UUID.Valid {
		return nil, errors.New("email verification wrong user.UUID")
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

	key := strings.Join([]string{session.RefreshTokenKey, refreshToken.UUID, refreshToken.Token}, ":")
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
	u.Email = types.NewNullString(req.Email)

	u, err := a.userRepository.FindByEmail(ctx, u)
	if err != nil {
		return nil, err
	}
	if !u.UUID.Valid {
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
	a.Cache().Set(fmt.Sprintf(PhoneRegistrationKey, phone), &code, 15*time.Minute)
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
	phone, err := validators.CheckPhoneFormat(req.Phone)
	if err != nil {
		return nil, err
	}
	code := 3455
	if !a.Config().SMSC.Dev {
		err = a.Cache().Get(fmt.Sprintf(PhoneRegistrationKey, phone), &code)
		if err != nil {
			return nil, err
		}
	}

	if code != req.Code {
		return nil, errors.New("phone code mismatch")
	}

	phoneEnt := types.NewNullString(phone)
	u := infoblog.User{
		Phone: phoneEnt,
	}
	u, err = a.userRepository.FindByPhone(ctx, u)
	if err != nil && err.Error() == "sql: no rows in result set" {
		u, err = createDefaultUser()
		if err != nil {
			return nil, err
		}
		u.Phone = phoneEnt

		err = a.userRepository.CreateUser(ctx, u)
		if err != nil {
			return nil, err
		}

		u, err = a.userRepository.Find(ctx, u)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	token, err := a.JWTKeys().GenerateAuthTokens(&u)
	if err != nil {
		return nil, err
	}

	return token, err
}

func (a *service) SocialCallback(ctx context.Context, state string) (string, error) {
	var err error
	u, ok := ctx.Value(types.User{}).(*infoblog.User)
	if !ok {
		return "", errors.New("type assertion to user err")
	}

	provider, ok := ctx.Value(Provider{}).(string)
	if !ok {
		return "", fmt.Errorf("wrong provider")
	}

	var user infoblog.User
	switch provider {
	case "facebook":
		user, err = a.userRepository.FindByFacebook(ctx, *u)
		if err != nil {
			user, err = createDefaultUser()
			if err != nil {
				return "", err
			}
			user.FacebookID = u.FacebookID
			user.Name = u.Name
			err = a.userRepository.CreateUser(ctx, user)
			if err != nil {
				return "", err
			}
		}
	case "google":
		user, err = a.userRepository.FindByGoogle(ctx, *u)
		if err != nil {
			user, err = createDefaultUser()
			if err != nil {
				return "", err
			}
			user.GoogleID = u.GoogleID
			user.Name = u.Name
			err = a.userRepository.CreateUser(ctx, user)
			if err != nil {
				return "", err
			}
		}
	default:
		return "", fmt.Errorf("unknown provider %s", provider)
	}

	a.Cache().Set(fmt.Sprintf(SocialsAuthKey, state), &user, 10*time.Minute)

	uri, err := url.Parse(a.Config().App.FrontEnd)
	if err != nil {
		return "", err
	}

	uri.Path = fmt.Sprintf("socials/%s", state)

	return uri.String(), nil
}

func (a *service) Oauth2Token(ctx context.Context, stateRequest request.StateRequest) (*request.AuthTokenData, error) {
	_ = ctx
	var u infoblog.User
	key := fmt.Sprintf(SocialsAuthKey, stateRequest.State)
	err := a.Cache().Get(key, &u)
	if err != nil {
		return nil, err
	}
	err = a.Cache().Del(key)
	if err != nil {
		a.Logger().Error("social auth key deletion", zap.Error(err))
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

func createDefaultUser() (infoblog.User, error) {
	return infoblog.User{
		UUID:        types.NewNullUUID(),
		Trial:       types.NewNullBool(true),
		NotifyEmail: types.NewNullBool(true),
		Language:    types.NewNullInt64(1),
	}, nil
}
