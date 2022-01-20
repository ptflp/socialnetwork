package controllers

import (
	"net/http"

	"gitlab.com/InfoBlogFriends/server/components"

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

func NewChatController(components components.Componenter, services *services.Services) *chatController {
	return &chatController{
		Decoder:   decoder.NewDecoder(),
		Responder: components.Responder(),
		logger:    components.Logger(),
		chats:     services.Chats,
	}
}

func (a *chatController) SendMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sendMessagePrivate request.SendMessageReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &sendMessagePrivate)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		chatData, err := a.chats.SendMessage(r.Context(), sendMessagePrivate)

		if err != nil {
			a.ErrorInternal(w, err)
			return
		}
		a.SendJSON(w, request.Response{
			Success: true,
			Data:    chatData,
		})
	}
}

func (a *chatController) GetMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getMessages request.GetMessagesReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &getMessages)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		messagesData, err := a.chats.GetMessages(r.Context(), getMessages)

		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    messagesData,
		})
	}
}

func (a *chatController) GetChats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getChatsReq request.GetChatsReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &getChatsReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		chatsData, err := a.chats.GetChats(r.Context(), getChatsReq)

		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    chatsData,
		})
	}
}

func (a *chatController) Info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getInfoReq request.GetInfoReq
		// r.PostForm is a map of our POST form values
		err := Decode(r, &getInfoReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		chatData, err := a.chats.Info(r.Context(), getInfoReq)

		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    chatData,
		})
	}
}
