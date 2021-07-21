package controllers

import (
	"net/http"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type usersController struct {
	*decoder.Decoder
	respond.Responder
	user   *services.User
	logger *zap.Logger
}

func NewUsersController(responder respond.Responder, user *services.User, logger *zap.Logger) *usersController {
	return &usersController{
		Decoder:   decoder.NewDecoder(),
		Responder: responder,
		user:      user,
		logger:    logger,
	}
}

func (u *usersController) Subscribe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := extractUser(r)
		if err != nil {
			u.ErrorBadRequest(w, err)
			return
		}

		var usersSubscribeReq request.UserIDRequest

		// r.PostForm is u map of our POST form values
		err = u.Decode(r.Body, &usersSubscribeReq)

		if err != nil {
			u.ErrorBadRequest(w, err)
			return
		}

		err = u.user.Subscribe(r.Context(), user, usersSubscribeReq)

		if err != nil {
			u.ErrorBadRequest(w, err)
			return
		}

		u.SendJSON(w, request.Response{
			Success: true,
		})
	}
}

func (u *usersController) Unsubscribe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := extractUser(r)
		if err != nil {
			u.ErrorBadRequest(w, err)
			return
		}

		var usersSubscribeReq request.UserIDRequest

		// r.PostForm is u map of our POST form values
		err = u.Decode(r.Body, &usersSubscribeReq)

		if err != nil {
			u.ErrorBadRequest(w, err)
			return
		}

		err = u.user.Unsubscribe(r.Context(), user, usersSubscribeReq)

		if err != nil {
			u.ErrorBadRequest(w, err)
			return
		}

		u.SendJSON(w, request.Response{
			Success: true,
		})
	}
}

func (u *usersController) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		usersData, err := u.user.List(r.Context())

		if err != nil {
			u.ErrorBadRequest(w, err)
			return
		}

		u.SendJSON(w, request.Response{
			Success: true,
			Data: struct {
				Users []request.UserData `json:"users"`
			}{
				Users: usersData,
			},
		})
	}
}

func (u *usersController) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var userIDNickReq request.UserIDNickRequest

		// r.PostForm is u map of our POST form values
		err := u.Decode(r.Body, &userIDNickReq)

		userData, err := u.user.Get(r.Context(), userIDNickReq)

		if err != nil {
			u.ErrorBadRequest(w, err)
			return
		}

		u.SendJSON(w, request.Response{
			Success: true,
			Data: struct {
				Users request.UserData `json:"user"`
			}{
				Users: userData,
			},
		})
	}
}
