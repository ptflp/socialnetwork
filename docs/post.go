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

// swagger:route POST /posts/feed/recent posts postsListRequest
// Получение ленты последних постов.
// security:
//   - Bearer: []
// responses:
//   200: postsListResponse

// swagger:response postsListResponse
type postsListResponse struct {
	// in:body
	Body request.PostsFeedResponse
}

// swagger:parameters postsListRequest
type postsListParams struct {
	// in:body
	FormData request.PostsFeedReq
}
