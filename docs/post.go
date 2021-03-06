package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /posts/create posts postsCreateRequest
// Создание поста.
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
	Body request.PostCreateReq
}

// swagger:route POST /posts/file/upload posts postFileUploadRequest
// Загрузка файла поста.
// security:
//   - Bearer: []
// responses:
//   200: postFileUploadResponse

// swagger:response postFileUploadResponse
type postFileUploadResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters postFileUploadRequest
type postFileUploadParams struct {
	// in: formData
	// swagger:file
	File interface{} `json:"file"`
}

// swagger:route POST /posts/file/upload/video posts postVideoUploadRequest
// Загрузка видео поста.
// security:
//   - Bearer: []
// responses:
//   200: postVideoUploadResponse

// swagger:response postVideoUploadResponse
type postVideoUploadResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters postVideoUploadRequest
type postVideoUploadParams struct {
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
	FormData request.LimitOffsetReq
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
	Body request.LimitOffsetReq
}

// swagger:route POST /posts/feed/user posts postsFeedUserRequest
// Получение ленты постов пользователя.
// security:
//   - Bearer: []
// responses:
//   200: postsFeedUserResponse

// swagger:response postsFeedUserResponse
type postsFeedUserResponse struct {
	// in:body
	Body request.PostsFeedResponse
}

// swagger:parameters postsFeedUserRequest
type postsFeedUserParams struct {
	// in:body
	Body request.PostsFeedUserReq
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
	Body request.LikeReq
}

// swagger:route POST /posts/get posts postGetRequest
// Получение поста по post_id.
// security:
//   - Bearer: []
// responses:
//   200: postGetResponse

// swagger:response postGetResponse
type postGetResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters postGetRequest
type postGetParams struct {
	// in:body
	Body request.PostUUIDReq
}

// swagger:route POST /posts/delete posts postDeleteRequest
// Удаление поста по post_id.
// security:
//   - Bearer: []
// responses:
//   200: postDeleteResponse

// swagger:response postDeleteResponse
type postDeleteResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters postDeleteRequest
type postDeleteParams struct {
	// in:body
	Body request.PostUUIDReq
}

// swagger:route POST /posts/update posts postUpdateRequest
// Обновление данных поста по post_id.
// security:
//   - Bearer: []
// responses:
//   200: postUpdateResponse

// swagger:response postUpdateResponse
type postUpdateResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters postUpdateRequest
type postUpdateParams struct {
	// in:body
	Body request.PostUpdateReq
}

// swagger:route POST /posts/feed/subscribed posts postsFeedSubscribedRequest
// Получение ленты своих подписок.
// security:
//   - Bearer: []
// responses:
//   200: postsFeedSubscribedResponse

// swagger:response postsFeedSubscribedResponse
type postsFeedSubscribedResponse struct {
	// in:body
	Body request.PostsFeedResponse
}

// swagger:parameters postsFeedSubscribedRequest
type postsFeedSubscribedParams struct {
	// in:body
	Body request.LimitOffsetReq
}

// swagger:route POST /posts/feed/recommends posts postsFeedRecommendsRequest
// Получение ленты рекомендации.
// security:
//   - Bearer: []
// responses:
//   200: postsFeedRecommendsResponse

// swagger:response postsFeedRecommendsResponse
type postsFeedRecommendsResponse struct {
	// in:body
	Body request.PostsFeedResponse
}

// swagger:parameters postsFeedRecommendsRequest
type postsFeedRecommendsParams struct {
	// in:body
	Body request.LimitOffsetReq
}

// swagger:route POST /expose/feed/recommends posts publicFeedRequest
// Получение ленты рекомендации.
// responses:
//   200: publicFeedResponse

// swagger:response publicFeedResponse
type publicFeedResponse struct {
	// in:body
	Body request.PostsFeedResponse
}

// swagger:parameters publicFeedRequest
type publicFeedParams struct {
	// in:body
	Body request.LimitOffsetReq
}

// swagger:route POST /posts/comments/create posts CommentCreateRequest
// Создание комментария к посту.
// responses:
//   200: CommentCreateResponse

// swagger:response CommentCreateResponse
type CommentCreateResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters CommentCreateRequest
type CommentCreateParams struct {
	// in:body
	Body request.CommentCreateReq
}

// swagger:route POST /posts/comments/reply posts CommentReplyRequest
// Создание комментария к посту.
// responses:
//   200: CommentReplyResponse

// swagger:response CommentReplyResponse
type CommentReplyResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters CommentReplyRequest
type CommentReplyParams struct {
	// in:body
	Body request.CommentReplyReq
}

// swagger:route POST /posts/comments/get posts GetCommentsRequest
// Получение комментариев поста.
// responses:
//   200: GetCommentsResponse

// swagger:response GetCommentsResponse
type GetCommentsResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters GetCommentsRequest
type GetCommentsParams struct {
	// in:body
	Body request.PostUUIDReq
}
