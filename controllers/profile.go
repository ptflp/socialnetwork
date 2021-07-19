package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

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
		u, ok := ctx.Value("user").(*infoblog.User)
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

		*u, err = a.user.UpdateProfile(ctx, profileUpdateReq, *u)
		if err != nil {
			a.logger.Error("user service update profile", zap.Error(err))
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Msg:     "Данные профиля обновлены",
			Data: struct {
				User infoblog.User
			}{
				*u,
			},
		})
	}
}

func (a *profileController) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, ok := ctx.Value("user").(*infoblog.User)
		if !ok {
			a.ErrorInternal(w, errors.New("type assertion to user err"))
			return
		}
		user, err := a.user.GetProfile(ctx, *u)
		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Msg:     "",
			Data: struct {
				User infoblog.User `json:"user"`
			}{
				user,
			},
		})
	}
}

func (a *profileController) SetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, ok := ctx.Value("user").(*infoblog.User)
		if !ok {
			a.ErrorInternal(w, errors.New("type assertion to user err"))
			return
		}

		var setPasswordReq request.SetPasswordReq
		err := json.NewDecoder(r.Body).Decode(&setPasswordReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		err = a.user.SetPassword(ctx, *u)
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
