package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /user/subscribe user userSubscribeRequest
// Обновление профиля.
// security:
//   - Bearer: []
// responses:
//   200: userSubscribeResponse

// swagger:response userSubscribeResponse
type userSubscribeResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters userSubscribeRequest
type userSubscribeParams struct {
	// in:body
	Body request.UserSubscribeRequest
}
