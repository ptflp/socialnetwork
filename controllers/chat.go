package controllers

import (
	"net/http"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type chatsController struct {
	*decoder.Decoder
	respond.Responder
	chat   *services.Chat
	logger *zap.Logger
}

func NewChatsController(responder respond.Responder, chat *services.Chat, logger *zap.Logger) *chatsController {
	return &chatsController{
		Decoder:   decoder.NewDecoder(),
		Responder: responder,
		chat:      chat,
		logger:    logger,
	}
}

func (c *chatsController) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		chatsData, err := c.chat.List(r.Context())

		if err != nil {
			c.ErrorBadRequest(w, err)
			return
		}

		c.SendJSON(w, request.Response{
			Success: true,
			Data: struct {
				Chats []request.ChatData `json:"chats"`
			}{
				Chats: chatsData,
			},
		})
	}
}

func (c *chatsController) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var chatIDReq request.ChatIDRequest

		err := c.Decode(r.Body, &chatIDReq)
		if err != nil {
			c.ErrorBadRequest(w, err)
			return
		}

		chatData, err := c.chat.Get(r.Context(), chatIDReq)

		if err != nil {
			c.ErrorBadRequest(w, err)
			return
		}

		c.SendJSON(w, request.Response{
			Success: true,
			Data: struct {
				Chats request.ChatData `json:"chat"`
			}{
				Chats: chatData,
			},
		})
	}
}
