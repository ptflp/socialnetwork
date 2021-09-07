package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /chat/message/send chatMessage chatMessageRequest
// Signal.
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
	Body request.SendMessage
}
