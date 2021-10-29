package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /notification/my notification notificationRequest
// Получение всех своих нотификаций.
// responses:
//   200: notificationResponse

// swagger:response notificationResponse
type notificationResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters notificationRequest
type notificationParams struct {
	// in:body
	Body request.LimitOffsetReq
}

// swagger:route POST /notification/shown notification notificationShownRequest
// Установка флага просмотренно по id нотификации/эвента.
// responses:
//   200: notificationShownResponse

// swagger:response notificationShownResponse
type notificationShownResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters notificationShownRequest
type notificationShownParams struct {
	// in:body
	Body []request.UUIDReq
}
