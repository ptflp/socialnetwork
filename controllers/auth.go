package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/request"

	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type authController struct {
	respond.Responder
	authService infoblog.AuthService
	logger      *zap.Logger
}

func NewAuth(responder respond.Responder, authService infoblog.AuthService, logger *zap.Logger) *authController {
	return &authController{
		Responder:   responder,
		authService: authService,
		logger:      logger,
	}
}

func (a *authController) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshTokenReq request.RefreshTokenRequest
		err := json.NewDecoder(r.Body).Decode(&refreshTokenReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		token, err := a.authService.RefreshToken(r.Context(), &refreshTokenReq)
		if err != nil {
			a.ErrorForbidden(w, err)
			return
		}
		a.SendJSON(w, request.AuthTokenResponse{
			Success: false,
			Msg:     "",
			Data:    *token,
		})
	}
}

func (a *authController) EmailActivation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var emailActivationReq request.EmailActivationRequest
		err := json.NewDecoder(r.Body).Decode(&emailActivationReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		if err = a.authService.EmailActivation(r.Context(), &emailActivationReq); err != nil {
			a.SendJSON(w, request.Response{
				Success: false,
				Msg:     fmt.Sprintf("Ошибка отправки почты: %s", err),
				Data:    nil,
			})
			return
		}
		a.SendJSON(w, request.Response{
			Success: true,
			Msg:     fmt.Sprintf("Ссылка активации отправлена на почту %s", emailActivationReq.Email),
			Data:    nil,
		})
	}
}

func (a *authController) EmailVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var emailVerificationReq request.EmailVerificationRequest
		err := json.NewDecoder(r.Body).Decode(&emailVerificationReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		token, err := a.authService.EmailVerification(r.Context(), &emailVerificationReq)
		if err != nil {
			a.SendJSON(w, request.Response{
				Success: false,
				Msg:     fmt.Sprintf("email verification err: %s", err),
				Data:    nil,
			})
			return
		}

		a.SendJSON(w, request.AuthTokenResponse{
			Success: true,
			Msg:     "Ваша почта успешно активирована",
			Data:    *token,
		})
	}
}

func (a *authController) EmailLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var emailLoginReq request.EmailLoginRequest
		err := json.NewDecoder(r.Body).Decode(&emailLoginReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		token, err := a.authService.EmailLogin(r.Context(), &emailLoginReq)
		if err != nil {
			a.SendJSON(w, request.Response{
				Success: false,
				Msg:     fmt.Sprintf("email login err: %s", err),
				Data:    nil,
			})
			return
		}

		a.SendJSON(w, request.AuthTokenResponse{
			Success: true,
			Msg:     "",
			Data:    *token,
		})
	}
}

func (a *authController) SendCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sendCodeReq request.PhoneCodeRequest
		err := json.NewDecoder(r.Body).Decode(&sendCodeReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		if a.authService.SendCode(r.Context(), &sendCodeReq) {
			a.SendJSON(w, request.Response{
				Success: true,
				Msg:     "СМС код оптравлен успешно",
				Data:    nil,
			})
			return
		}
		a.SendJSON(w, request.Response{
			Success: false,
			Msg:     "Ошибка отправки кода",
			Data:    nil,
		})
	}
}

func (a *authController) CheckCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var checkCodeReq request.CheckCodeRequest
		err := json.NewDecoder(r.Body).Decode(&checkCodeReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		token, err := a.authService.CheckCode(r.Context(), &checkCodeReq)
		if err != nil {
			a.SendJSON(w, request.Response{
				Success: false,
				Msg:     fmt.Sprintf("Ошибка проверки кода: %s", err),
				Data:    nil,
			})
			return
		}
		a.SendJSON(w, request.AuthTokenResponse{
			Success: true,
			Msg:     "",
			Data:    *token,
		})
	}
}

func (a *authController) Oauth2Callback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		url, _ := a.authService.SocialCallback(r.Context(), state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (a *authController) Oauth2State() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var stateReq request.StateRequest
		err := json.NewDecoder(r.Body).Decode(&stateReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
		}
		token, err := a.authService.Oauth2Token(r.Context(), stateReq)
		if err != nil {
			a.SendJSON(w, request.Response{
				Success: false,
				Msg:     fmt.Sprintf("Ошибка проверки кода: %s", err),
				Data:    nil,
			})
			return
		}
		a.SendJSON(w, request.AuthTokenResponse{
			Success: true,
			Msg:     "",
			Data:    *token,
		})
	}
}
