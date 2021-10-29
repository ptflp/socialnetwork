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

type notificationController struct {
	*decoder.Decoder
	respond.Responder
	logger        *zap.Logger
	notifications *services.Event
}

func NewNotificationController(components components.Componenter, services *services.Services) *notificationController {
	return &notificationController{
		Decoder:       decoder.NewDecoder(),
		Responder:     components.Responder(),
		logger:        components.Logger(),
		notifications: services.Event,
	}
}

func (a *notificationController) GetMy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var limitOffsetReq request.LimitOffsetReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &limitOffsetReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		notificationData, err := a.notifications.GetMy(r.Context(), limitOffsetReq)

		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    notificationData,
		})
	}
}

func (a *notificationController) Shown() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var uuidReqs []request.UUIDReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &uuidReqs)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		err = a.notifications.Shown(r.Context(), uuidReqs)

		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
		})
	}
}
