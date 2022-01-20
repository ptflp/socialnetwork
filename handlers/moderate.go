package handlers

import (
	"net/http"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type moderateController struct {
	*decoder.Decoder
	respond.Responder
	user      *services.User
	file      *services.File
	post      *services.Post
	logger    *zap.Logger
	comments  *services.Comments
	moderates *services.Moderates
}

func NewModerateController(responder respond.Responder, services *services.Services, logger *zap.Logger) *moderateController {
	return &moderateController{
		Decoder:   decoder.NewDecoder(),
		Responder: responder,
		user:      services.User,
		file:      services.File,
		post:      services.Post,
		comments:  services.Comments,
		logger:    logger,
		moderates: services.Moderates,
	}
}

func (a *moderateController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createModerate request.ModerateCreateReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &createModerate)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		err = a.moderates.CreateModerate(r.Context(), createModerate)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
		})
	}
}

func (a *moderateController) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var moderateUUID request.PostUUIDReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &moderateUUID)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		moderate, err := a.moderates.Get(r.Context(), moderateUUID)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    moderate,
		})
	}
}

func (a *moderateController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var moderateUUID request.ModerateUpdateStatusReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &moderateUUID)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		moderate, err := a.moderates.Update(r.Context(), moderateUUID)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    moderate,
		})
	}
}

func (a *moderateController) GetModerates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var limitOffsetReq request.LimitOffsetReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &limitOffsetReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		moderate, err := a.moderates.GetModerates(r.Context(), limitOffsetReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    moderate,
		})
	}
}

func (a *moderateController) UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(100 << 20)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		file, fHeader, err := r.FormFile("file")
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		defer file.Close()

		formFile := services.FormFile{
			File:       file,
			FileHeader: fHeader,
		}

		fileData, err := a.moderates.SaveFile(r.Context(), formFile)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    fileData,
		})
	}
}
