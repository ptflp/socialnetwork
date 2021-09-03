package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /restricted/signal signal signalRequest
// Signal.
// security:
//   - Bearer: []
// responses:
//   200: signalResponse

// swagger:response signalResponse
type signalResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters signalRequest
type signalParams struct {
	// in:body
	Body interface{}
}
