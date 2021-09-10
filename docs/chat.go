package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /chat/message/send chat chatMessageRequest
// Отправить сообщение.
// security:
//   - Bearer: []
// responses:
//   200: chatMessageResponse

// swagger:response chatMessageResponse
type chatMessageResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters chatMessageRequest
type chatMessageParams struct {
	// in:body
	Body request.SendMessageReq
}

// swagger:route POST /chat/info chat chatInfoRequest
// Получить информацию о чате.
// security:
//   - Bearer: []
// responses:
//   200: chatInfoResponse

// swagger:response chatInfoResponse
type chatInfoResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters chatInfoRequest
type chatInfoParams struct {
	// in:body
	Body request.GetInfoReq
}

// swagger:route POST /chat/get/messages chat chatGetMessagesRequest
// Получить информацию о чате.
// security:
//   - Bearer: []
// responses:
//   200: chatGetMessagesResponse

// swagger:response chatGetMessagesResponse
type chatGetMessagesResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters chatGetMessagesRequest
type chatGetMessagesParams struct {
	// in:body
	Body request.GetMessagesReq
}
