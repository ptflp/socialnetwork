package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /people/subscribe people userSubscribeRequest
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
	Body request.UserIDRequest
}

// swagger:route POST /people/unsubscribe people userUnsubscribeRequest
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
	Body request.UserIDRequest
}

// swagger:route GET /people/list people userListRequest
// Лист пользователей.
// security:
//   - Bearer: []
// responses:
//   200: userListResponse

// swagger:response userListResponse
type userListResponse struct {
	// in:body
	Body request.Response
}

// swagger:route POST /people/get people peopleIDRequest
// Получить пользователя по никнейму или айди.
// security:
//   - Bearer: []
// responses:
//   200: peopleIDResponse

// swagger:response peopleIDResponse
type peopleIDResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters peopleIDRequest
type peopleIDParams struct {
	// in:body
	Body request.UserIDNickRequest
}

// swagger:route GET /people/list/recommends people userRecommendsRequest
// Лист пользователей.
// security:
//   - Bearer: []
// responses:
//   200: userRecommendsResponse

// swagger:response userRecommendsResponse
type userRecommendsResponse struct {
	// in:body
	Body request.Response
}

// swagger:route GET /people/list/subscribed people userSubscribedListRequest
// Лист пользователей.
// security:
//   - Bearer: []
// responses:
//   200: userSubscribedListResponse

// swagger:response userSubscribedListResponse
type userSubscribedListResponse struct {
	// in:body
	Body request.Response
}

// swagger:route POST /people/autocomplete people peopleAutocompleteRequest
// Автозаполнение по никнейму.
// security:
//   - Bearer: []
// responses:
//   200: peopleAutocompleteResponse

// swagger:response peopleAutocompleteResponse
type peopleAutocompleteResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters peopleAutocompleteRequest
type peopleAutocompleteParams struct {
	// in:body
	Body request.UserNicknameRequest
}
