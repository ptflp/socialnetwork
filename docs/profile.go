package docs

import (
	"gitlab.com/InfoBlogFriends/server/request"
)

// swagger:route POST /profile/update profile profileUpdateRequest
// Обновление профиля.
// security:
//   - Bearer: []
// responses:
//   200: profileUpdateResponse

// swagger:response profileUpdateResponse
type profileUpdateResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters profileUpdateRequest
type profileUpdateParams struct {
	// in:body
	Body request.ProfileUpdateReq
}

// swagger:route POST /profile/set/password profile profileSetPassword
// Изменение пароля.
// security:
//   - Bearer: []
// responses:
//   200: profileSetPasswordResponse

// swagger:response profileSetPasswordResponse
type profileSetPasswordResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters profileSetPassword
type profileSetPasswordParams struct {
	// in:body
	Body request.SetPasswordReq
}

// swagger:route GET /profile/get profile getProfileRequest
// Получение профиля.
// security:
//   - Bearer: []
// responses:
//   200: getProfileResponse

// swagger:response getProfileResponse
type getProfileResponse struct {
	// in:body
	Body request.Response
}

// swagger:route POST /profile/upload/avatar profile avatarUploadRequest
// Загрузка аватарки пользователя.
// security:
//   - Bearer: []
// responses:
//   200: avatarUploadResponse

// swagger:response avatarUploadResponse
type avatarUploadResponse struct {
	// in:body
	Body request.Response
}

// swagger:parameters avatarUploadRequest
type avatarUploadParams struct {
	// in: formData
	// swagger:file
	File interface{} `json:"file"`
}
