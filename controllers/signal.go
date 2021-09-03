package controllers

import (
	"net/http"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type signalController struct {
	*decoder.Decoder
	respond.Responder
	user      *services.User
	file      *services.File
	post      *services.Post
	logger    *zap.Logger
	comments  *services.Comments
	moderates *services.Moderates
}

func NewSignalController(responder respond.Responder, services *services.Services, logger *zap.Logger) *signalController {
	return &signalController{
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

func (a *signalController) Signal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.SendJSON(w, request.Response{
			Success: true,
		})
	}
}
