package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/types"

	"gitlab.com/InfoBlogFriends/server/request"

	infoblog "gitlab.com/InfoBlogFriends/server"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type profileController struct {
	respond.Responder
	user   *services.User
	logger *zap.Logger
}

func NewProfileController(responder respond.Responder, user *services.User, logger *zap.Logger) *profileController {
	return &profileController{
		Responder: responder,
		user:      user,
		logger:    logger,
	}
}

func (a *profileController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, ok := ctx.Value(types.User{}).(*infoblog.User)
		if !ok {
			a.ErrorInternal(w, errors.New("type assertion to user err"))
			return
		}

		var profileUpdateReq request.ProfileUpdateReq
		err := json.NewDecoder(r.Body).Decode(&profileUpdateReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		userData, err := a.user.UpdateProfile(ctx, profileUpdateReq, *u)
		if err != nil {
			a.logger.Error("user service update profile", zap.Error(err))
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Msg:     "Данные профиля обновлены",
			Data: struct {
				User request.UserData `json:"user"`
			}{
				userData,
			},
		})
	}
}

func (a *profileController) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := a.user.GetProfile(r.Context())
		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Msg:     "",
			Data: struct {
				User request.UserData `json:"user"`
			}{
				user,
			},
		})
	}
}

func (a *profileController) SetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var setPasswordReq request.SetPasswordReq
		err := json.NewDecoder(r.Body).Decode(&setPasswordReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		err = a.user.SetPassword(r.Context(), setPasswordReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Msg:     "пароль успешно изменен",
			Data:    nil,
		})
	}
}

func (a *profileController) UploadAvatar() http.HandlerFunc {
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

		userData, err := a.user.SaveAvatar(r.Context(), formFile)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    userData,
		})
	}
}
