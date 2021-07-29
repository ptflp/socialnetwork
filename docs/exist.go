package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /exist/email exist existEmailRequest
// Проверка на существование почты.
// responses:
//   200: existEmailResponse

// swagger:response existEmailResponse
type existEmailResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters existEmailRequest
type existEmailParams struct {
	// in:body
	Body request.EmailRequest
}

// swagger:route POST /exist/nickname exist existNicknameRequest
// Проверка на существование никнейма.
// responses:
//   200: existNicknameResponse

// swagger:response existNicknameResponse
type existNicknameResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters existNicknameRequest
type existNicknameParams struct {
	// in:body
	Body request.NicknameRequest
}
