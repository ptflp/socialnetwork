package docs

import "gitlab.com/InfoBlogFriends/server/request"

// swagger:route POST /moderate/create moderate moderateCreateRequest
// Создание поста.
// security:
//   - Bearer: []
// responses:
//   200: moderateCreateResponse

// swagger:response moderateCreateResponse
type moderateCreateResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters moderateCreateRequest
// Создание записи модерации.
type moderateCreateParams struct {
	// in:body
	Body request.ModerateCreateReq
}

// swagger:route POST /moderate/file/upload moderate moderateFileUploadRequest
// Загрузка файла поста.
// security:
//   - Bearer: []
// responses:
//   200: moderateFileUploadResponse

// swagger:response moderateFileUploadResponse
type moderateFileUploadResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters moderateFileUploadRequest
type moderateFileUploadParams struct {
	// in: formData
	// swagger:file
	File interface{} `json:"file"`
}

// swagger:route POST /moderate/get moderate moderateGetRequest
// Создание поста.
// security:
//   - Bearer: []
// responses:
//   200: moderateGetResponse

// swagger:response moderateGetResponse
type moderateGetResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters moderateGetRequest
// Создание записи модерации.
type moderateGetParams struct {
	// in:body
	Body request.UUIDReq
}

// swagger:route POST /moderate/get/all moderate moderateGetAllRequest
// Создание поста.
// security:
//   - Bearer: []
// responses:
//   200: moderateGetAllResponse

// swagger:response moderateGetAllResponse
type moderateGetAllResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters moderateGetAllRequest
// Создание записи модерации.
type moderateGetAllParams struct {
	// in:body
	Body request.LimitOffsetReq
}
