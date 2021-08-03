package docs

import (
	"gitlab.com/InfoBlogFriends/server/request"
)

// swagger:route POST /auth/code auth sendCodeRequest
// Отправка смс кода.
// responses:
//   200: sendCodeResponse

// swagger:response sendCodeResponse
type sendCodeResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters sendCodeRequest
type sendCodeParams struct {
	// in:body
	Body request.PhoneCodeRequest
}

// swagger:route POST /auth/checkcode auth checkCodeRequest
// Проверка смс кода.
// responses:
//   200: checkCodeResponse

// swagger:response checkCodeResponse
type checkCodeResponse struct {
	// in:body
	Body request.AuthTokenData
}

// swagger:parameters checkCodeRequest
type checkCodeParams struct {
	// in:body
	Body request.CheckCodeRequest
}

// swagger:route POST /auth/email/registration auth EmailActivationRequest
// Отправка ссылки активации на почту.
// responses:
//   200: emailRegistrationResponse

// swagger:response emailRegistrationResponse
type emailRegistrationResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters EmailActivationRequest
type emailActivationParams struct {
	// in:body
	Body request.EmailActivationRequest
}

// swagger:route POST /auth/email/verification auth EmailVerificationRequest
// Подтверждение почты, авторизация.
// responses:
//   200: EmailVerificationResponse

// swagger:response EmailVerificationResponse
type emailVerificationResponse struct {
	// in:body
	Body request.AuthTokenResponse
}

// swagger:parameters EmailVerificationRequest
type emailVerificationParams struct {
	// in:body
	Body request.EmailVerificationRequest
}

// swagger:route POST /auth/email/login auth EmailLoginRequest
// Авторизация пользователя по емейл + пароль.
// responses:
//   200: EmailLoginResponse

// swagger:response EmailLoginResponse
type emailLoginResponse struct {
	// in:body
	Body request.AuthTokenResponse
}

// swagger:parameters EmailLoginRequest
type emailLoginParams struct {
	// in:body
	Body request.EmailLoginRequest
}

// swagger:route POST /auth/token/refresh auth RefreshTokenRequest
// Обновление токена.
// responses:
//   200: RefreshResponse

// swagger:response RefreshResponse
type refTokenResponse struct {
	// in:body
	Body request.AuthTokenResponse
}

// swagger:parameters RefreshTokenRequest
type refTokenParams struct {
	// in:body
	Body request.RefreshTokenRequest
}

// swagger:route POST /auth/oauth2/state auth Oauth2StateRequest
// Авторизация с помощью state oauth2.
// responses:
//   200: OauthStateResponse

// swagger:response OauthStateResponse
type oauthStateResponse struct {
	// in:body
	Body request.AuthTokenResponse
}

// swagger:parameters Oauth2StateRequest
type oauthStateParams struct {
	// in:body
	Body request.StateRequest
}

// swagger:route GET /auth/oauth2/facebook/login auth FacebookLoginRequest
// Авторизация с помощью фэйсбук.
// responses:
//   200: FacebookLoginResponse

// swagger:response FacebookLoginResponse
type facebookLoginResponse struct {
	// in:body
	Body request.Response
}

// swagger:route GET /auth/oauth2/google/login auth GoogleLoginRequest
// Авторизация с помощью гугл.
// responses:
//   200: GoogleLoginResponse

// swagger:response GoogleLoginResponse
type googleLoginResponse struct {
	// in:body
	Body request.Response
}
