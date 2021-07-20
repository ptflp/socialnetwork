package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /user/subscribe user userSubscribeRequest
// Подписаться на пользователя.
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
	Body request.UserSubscriberRequest
}

// swagger:route POST /user/unsubscribe user userUnsubscribeRequest
// Отписаться от пользователя.
// security:
//   - Bearer: []
// responses:
//   200: userUnsubscribeResponse

// swagger:response userUnsubscribeResponse
type userUnsubscribeResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters userUnsubscribeRequest
type userUnsubscribeParams struct {
	// in:body
	Body request.UserSubscriberRequest
}

// swagger:route GET /user/list user userListRequest
// Отписаться от пользователя.
// security:
//   - Bearer: []
// responses:
//   200: userListResponse

// swagger:response userListResponse
type userListResponse struct {
	// in:body
	Body request.Response
}
