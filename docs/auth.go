package docs

import (
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/server"
)

// swagger:route POST /auth/code auth sendCodeRequest
// Отправка смс кода.
// responses:
//   200: sendCodeResponse

// swagger:response sendCodeResponse
type sendCodeResponse struct {
	// in:body
	Body server.Response
}

// swagger:parameters sendCodeRequest
type sendCodeParams struct {
	// in:body
	Body infoblog.PhoneCodeRequest
}

// swagger:route POST /auth/checkcode auth checkCodeRequest
// Проверка смс кода
// responses:
//   200: checkCodeResponse

// swagger:response checkCodeResponse
type checkCodeResponse struct {
	// in:body
	Body server.Response
}

// swagger:parameters checkCodeRequest
type checkCodeParams struct {
	// in:body
	Body infoblog.CheckCodeRequest
}
