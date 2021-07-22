package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /exist/email exist existEmailRequest
// Проверка на существование почты.
// security:
//   - Bearer: []
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
