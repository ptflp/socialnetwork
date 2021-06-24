package server

import (
	"net/http"

	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type authHandler struct {
	respond.Responder
	authService infoblog.AuthService
	logger      *zap.Logger
}

func NewAuthHandler(responder respond.Responder, authService infoblog.AuthService, logger *zap.Logger) (*authHandler, error) {
	return &authHandler{
		Responder:   responder,
		authService: authService,
		logger:      logger,
	}, nil
}

func (a *authHandler) SendCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.Responder.SendJSON(w, struct {
		}{})
	}
}
