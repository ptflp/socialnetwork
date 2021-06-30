package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	infoblog "gitlab.com/InfoBlogFriends/server"

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

type ProfileUpdateReq struct {
	Phone      *string `json:"phone"`
	Email      *string `json:"email"`
	Password   *string `json:"password"`
	Active     *string `json:"active"`
	Name       *string `json:"name"`
	SecondName *string `json:"second_name"`
}

func (a *profileHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, ok := ctx.Value("user").(*infoblog.User)
		if !ok {
			a.ErrorInternal(w, errors.New("type assertion to user err"))
			return
		}
		var sendCodeReq ProfileUpdateReq
		err := json.NewDecoder(r.Body).Decode(&sendCodeReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		a.SendJSON(w, u)
	}
}
