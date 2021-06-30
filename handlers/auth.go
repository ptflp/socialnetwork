package handlers

import (
	"encoding/json"
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

type PhoneCodeRequest struct {
	Phone string `json:"phone"`
}

func (a *authHandler) SendCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sendCodeReq PhoneCodeRequest
		err := json.NewDecoder(r.Body).Decode(&sendCodeReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		if a.authService.SendCode(r.Context(), &sendCodeReq) {
			a.Responder.SendJSON(w, Response{
				Success: true,
				Msg:     "СМС код оптравлен успешно",
				Data:    nil,
			})
			return
		}
		a.Responder.SendJSON(w, Response{
			Success: false,
			Msg:     "Ошибка отправки кода",
			Data:    nil,
		})
	}
}

type CheckCodeRequest struct {
	Phone string
	Code  int
}

func (a *authHandler) CheckCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var checkCodeReq CheckCodeRequest
		err := json.NewDecoder(r.Body).Decode(&checkCodeReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		token, err := a.authService.CheckCode(r.Context(), &checkCodeReq)
		if err != nil {
			a.Responder.SendJSON(w, Response{
				Success: false,
				Msg:     "Ошибка проверки кода " + err.Error(),
				Data:    nil,
			})
			return
		}
		a.Responder.SendJSON(w, Response{
			Success: true,
			Msg:     "",
			Data: struct {
				Token string
			}{
				Token: token,
			},
		})
	}
}

type Response struct {
	Success bool
	Msg     string
	Data    interface{}
}
