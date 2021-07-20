package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /posts/create posts postsCreateRequest
// Создание поста с файлом.
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
// Здесь все поля multipart/form-data, go-swagger не поддерживает формы, можно отправить запрос Postman
type postsCreateParams struct {
	// in:body
	FormData request.PostCreateReq
	// in: formData
	// swagger:file
	File interface{} `json:"file"`
}

// swagger:route POST /posts/feed/recent posts postsFeedRequest
// Получение ленты последних постов.
// security:
//   - Bearer: []
// responses:
//   200: postsFeedResponse

// swagger:response postsFeedResponse
type postsFeedResponse struct {
	// in:body
	Body request.PostsFeedResponse
}

// swagger:parameters postsFeedRequest
type postsFeedParams struct {
	// in:body
	FormData request.PostsFeedReq
}

// swagger:route POST /posts/feed/my posts postsFeedMyRequest
// Получение ленты своих постов.
// security:
//   - Bearer: []
// responses:
//   200: postsFeedMyResponse

// swagger:response postsFeedMyResponse
type postsFeedMyResponse struct {
	// in:body
	Body request.PostsFeedResponse
}

// swagger:parameters postsFeedMyRequest
type postsFeedMyParams struct {
	// in:body
	FormData request.PostsFeedReq
}

// swagger:route POST /posts/like posts postLikeRequest
// Лайк поста.
// security:
//   - Bearer: []
// responses:
//   200: postLikeResponse

// swagger:response postLikeResponse
type postLikeResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters postLikeRequest
type postLikeParams struct {
	// in:body
	FormData request.PostLikeReq
}
