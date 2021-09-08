package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type signalController struct {
	*decoder.Decoder
	respond.Responder
	user      *services.User
	file      *services.File
	post      *services.Post
	logger    *zap.Logger
	comments  *services.Comments
	moderates *services.Moderates
}

func NewSignalController(responder respond.Responder, services *services.Services, logger *zap.Logger) *signalController {
	return &signalController{
		Decoder:   decoder.NewDecoder(),
		Responder: responder,
		user:      services.User,
		file:      services.File,
		post:      services.Post,
		comments:  services.Comments,
		logger:    logger,
		moderates: services.Moderates,
	}
}

func (a *signalController) Signal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := extractUser(r)
		if err != nil {
			a.ErrorForbidden(w, err)
			return
		}
		data := Payload{
			Touserid: user.UUID.String,
			Data:     Data{Somedata: "asfasfasfafa"},
		}
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			return
		}
		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", "http://137.184.77.68/signalling", body)
		if err != nil {
			return
		}
		req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTYzMDY3NTQ4NSwiZXhwIjoxNjMwNjc5MDg1fQ.4tu5zY5ui1r5goHIQHqNqPEg9-IvYUDxxH7OCxJGKzQ")
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		a.SendJSON(w, request.Response{
			Success: true,
		})
	}
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
	Touserid string `json:"toUserId"`
	Data     Data   `json:"data"`
}
type Data struct {
	Somedata string `json:"someData"`
}
