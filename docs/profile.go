package docs

import (
	"gitlab.com/InfoBlogFriends/server/handlers"
)

// swagger:route POST /profile/update profile profileUpdateRequest
// Отправка смс кода.
// responses:
//   200: profileUpdateResponse

// swagger:response profileUpdateResponse
type profileUpdateResponse struct {
	// in:body
	Body handlers.Response
}

// swagger:parameters profileUpdateRequest
type profileUpdateParams struct {
	// in:body
	Body handlers.ProfileUpdateReq
}
