package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/service"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type profileHandler struct {
	respond.Responder
	user   service.User
	logger *zap.Logger
}

func NewProfileHandler(responder respond.Responder, user service.User, logger *zap.Logger) (*profileHandler, error) {
	return &profileHandler{
		Responder: responder,
		user:      user,
		logger:    logger,
	}, nil
}

func (a *profileHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sendCodeReq PhoneCodeRequest
		err := json.NewDecoder(r.Body).Decode(&sendCodeReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

	}
}

func (a *profileHandler) CheckCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var checkCodeReq CheckCodeRequest
		err := json.NewDecoder(r.Body).Decode(&checkCodeReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
	}
}
