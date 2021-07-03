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
// Подтверждение почты, авторизация
// responses:
//   200: EmailVerificationResponse

// swagger:response EmailVerificationResponse
type emailVerificationResponse struct {
	// in:body
	Body request.AuthTokenData
}

// swagger:parameters EmailVerificationRequest
type emailVerificationParams struct {
	// in:body
	Body request.EmailVerificationRequest
}
