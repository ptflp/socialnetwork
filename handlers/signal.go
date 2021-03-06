package handlers

import (
	"bytes"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/types"

	sq "github.com/Masterminds/squirrel"
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type signalController struct {
	*decoder.Decoder
	respond.Responder
	chat     *services.Chats
	event    *services.Event
	logger   *zap.Logger
	services *services.Services
}

func NewSignalController(responder respond.Responder, services *services.Services, logger *zap.Logger) *signalController {
	return &signalController{
		Decoder:   decoder.NewDecoder(),
		Responder: responder,
		chat:      services.Chats,
		event:     services.Event,
		logger:    logger,
		services:  services,
	}
}

func (a *signalController) Signal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := extractUser(r)
		if err != nil {
			a.ErrorForbidden(w, err)
			return
		}
		var wssReq WssRequest
		err = a.Decoder.Decode(r.Body, &wssReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		data := Payload{
			ToUserID: user.UUID.String,
		}
		ctx := r.Context()

		switch wssReq.Action {
		case ActionSendChatMessage:
			messageData, err := a.chat.SendMessage(ctx, request.SendMessageReq{
				Message:  wssReq.Message,
				ChatUUID: wssReq.UUID,
			})
			if err != nil {
				a.ErrorInternal(w, err)
				return
			}
			messagesData, err := a.chat.GetMessages(ctx, request.GetMessagesReq{ChatUUID: wssReq.UUID})
			if err != nil {
				a.ErrorInternal(w, err)
				return
			}
			data.Data.Res = request.Response{
				Success: true,
				Data:    messagesData,
			}
			data.Data.Action = ActionSendChatMessage
			chat, err := a.chat.GetInfo(ctx, request.GetInfoReq{
				ChatUUID: &wssReq.UUID,
			})
			if err == nil {
				condition := infoblog.Condition{
					Equal: &sq.Eq{"chat_uuid": types.NewNullUUID(wssReq.UUID)},
				}
				cpr, _ := a.chat.GetParticipants(ctx, condition)
				toUserData := data
				for i := range cpr {
					if cpr[i].UserUUID.String != user.UUID.String {
						toUserData.ToUserID = cpr[i].UserUUID.String
						_ = a.sendSignalMessage(toUserData)
						_, _ = a.services.Event.CreateEvent(r.Context(), services.ActionNotifyChatMessages, messageData.UUID, user.UUID, chat.Participants[i].UUID)
					}
				}
			}
		case ActionGetChats:
			chats, err := a.chat.GetChats(ctx, request.GetChatsReq{UserUUID: user.UUID.String})
			if err != nil {
				a.ErrorInternal(w, err)
				return
			}
			data.Data.Res = chats
			data.Data.Action = ActionGetChats
		case ActionGetChatInfo:
			chatData, err := a.chat.GetInfoByUser(ctx, request.GetInfoReq{
				UserUUID: &wssReq.UUID,
			})
			if err != nil {
				a.ErrorInternal(w, err)
				return
			}
			data.Data.Res = chatData
			data.Data.Action = ActionGetChatInfo
		case ActionGetChatMessages:
			messages, err := a.chat.GetMessages(ctx, request.GetMessagesReq{ChatUUID: wssReq.UUID})
			if err != nil {
				a.ErrorInternal(w, err)
				return
			}
			data.Data.Res = messages
			data.Data.Action = ActionGetChatMessages
		}

		err = a.sendSignalMessage(data)
		if err != nil {
			a.SendJSON(w, request.Response{
				Success: false,
			})
			return
		}

		a.SendJSON(w, request.Response{
			Success: err == nil,
		})
	}
}

func (a *signalController) sendSignalMessage(data Payload) error {
	var body bytes.Buffer
	var err error
	err = a.Encode(&body, data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://wss.fanam.org/signalling", &body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTYzMDY3NTQ4NSwiZXhwIjoxNjMwNjc5MDg1fQ.4tu5zY5ui1r5goHIQHqNqPEg9-IvYUDxxH7OCxJGKzQ")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

// curl --location --request POST 'http://137.184.77.68/signalling' \
// --header 'Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTYzMDY3NTQ4NSwiZXhwIjoxNjMwNjc5MDg1fQ.4tu5zY5ui1r5goHIQHqNqPEg9-IvYUDxxH7OCxJGKzQ' \
// --header 'Content-Type: application/json' \
// --data-raw '{
//     "toUserId" : "f115dd85-0cbd-11ec-82aa-0242ac17000a",
//     "data" : {
//         "someData" : "someData"
//     }
// }'

type Payload struct {
	Action   int    `json:"action"`
	ToUserID string `json:"toUserId"`
	Data     Data   `json:"data"`
}

type Data struct {
	Res    interface{} `json:"res"`
	Action int         `json:"action"`
}

const (
	ActionSendChatMessage = iota + 1
	ActionGetChats
	ActionGetChatMessages
	ActionGetChatInfo
)

type WssRequest struct {
	Action  int    `json:"action"`
	Message string `json:"message"`
	UUID    string `json:"uuid"`
}
