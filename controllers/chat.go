package controllers

import (
	"net/http"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type chatController struct {
	*decoder.Decoder
	respond.Responder
	logger *zap.Logger
	chats  *services.Chats
}

func NewChatController(responder respond.Responder, services *services.Services, logger *zap.Logger) *chatController {
	return &chatController{
		Decoder:   decoder.NewDecoder(),
		Responder: responder,
		logger:    logger,
		chats:     services.Chats,
	}
}

func (a *chatController) SendMessagePrivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sendMessagePrivate request.SendMessage

		// r.PostForm is a map of our POST form values
		err := Decode(r, &sendMessagePrivate)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		err = a.chats.SendMessagePrivate(r.Context(), sendMessagePrivate)

		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
		})
	}
}
