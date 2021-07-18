package controllers

import (
	"net/http"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type usersController struct {
	respond.Responder
	user   *services.User
	logger *zap.Logger
}

func NewUsersController(responder respond.Responder, user *services.User, logger *zap.Logger) *usersController {
	return &usersController{
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

		var usersSubscribeReq request.UserSubscribeRequest

		// r.PostForm is u map of our POST form values
		err = decoder.Decode(&usersSubscribeReq, r.PostForm)

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
