package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /posts/create posts postsCreateRequest
// Обновление профиля.
// security:
//   - Bearer: []
// responses:
//   200: postsCreateResponse

// swagger:response postsCreateResponse
type postsCreateResponse struct {
	// in:body
	Body request.PostDataResponse
}

// swagger:parameters postsCreateRequest
type postsCreateParams struct {
	// in:body
	FormData request.PostCreateReq
	// in: formData
	// swagger:file
	File interface{} `json:"file"`
}
