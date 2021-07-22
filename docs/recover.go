package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /recover/password recover recoverPasswordRequest
// Запрос на восстановление пароля по почте или телефону.
// security:
//   - Bearer: []
// responses:
//   200: recoverPasswordResponse

// swagger:response recoverPasswordResponse
type recoverPasswordResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters recoverPasswordRequest
type recoverPasswordParams struct {
	// in:body
	Body request.PasswordRecoverRequest
}

// swagger:route POST /recover/check/phone recover checkPhoneCodeRequest
// Проверка кода смс при запросе восстановления пароля.
// security:
//   - Bearer: []
// responses:
//   200: checkPhoneCodeResponse

// swagger:response checkPhoneCodeResponse
type checkPhoneCodeResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters checkPhoneCodeRequest
type checkPhoneCodeParams struct {
	// in:body
	Body request.CheckPhoneCodeRequest
}

// swagger:route POST /recover/set/password recover passwordResetRequest
// Установка пароля по recover_id, указывается в редиректе с почты /profile/password/{hash}, либо при проверке смс кода.
// security:
//   - Bearer: []
// responses:
//   200: passwordResetResponse

// swagger:response passwordResetResponse
type passwordResetResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters passwordResetRequest
type passwordResetParams struct {
	// in:body
	Body request.PasswordResetRequest
}
