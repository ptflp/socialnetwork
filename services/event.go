package services

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/components"
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"

	"gitlab.com/InfoBlogFriends/server/types"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type Event struct {
	eventRep infoblog.EventRepository
	user     *User
	ctx      context.Context
	logger   *zap.Logger
	ch       chan infoblog.Event
	*decoder.Decoder
}

func NewEventService(ctx context.Context, cmps components.Componenter, reps infoblog.Repositories, services *Services) *Event {
	v := &Event{Decoder: decoder.NewDecoder(), eventRep: reps.Events, ctx: ctx, logger: cmps.Logger(), ch: make(chan infoblog.Event, NumJobs), user: services.User}
	go v.WorkerPool()
	return v
}

func (e *Event) CreateEvent(ctx context.Context, eventType int, foreignUUID, userUUID, toUserUUID types.NullUUID) (infoblog.Event, error) {
	event := infoblog.Event{
		UUID:        types.NewNullUUID(),
		Type:        types.NewNullInt64(int64(eventType)),
		ForeignUUID: foreignUUID,
		UserUUID:    userUUID,
		ToUser:      toUserUUID,
		Notified:    types.NewNullBool(false),
		Active:      types.NewNullBool(true),
	}

	err := e.eventRep.Create(ctx, event)

	return event, err
}

func (e *Event) Shown(ctx context.Context, req []request.UUIDReq) error {
	user, err := extractUser(ctx)
	if err != nil {
		return err
	}
	uuids := make([]interface{}, 0, len(req))
	for i := range req {
		uuids = append(uuids, types.NewNullUUID(req[i].UUID))
	}

	condition := infoblog.Condition{
		Equal: &sq.Eq{"to_user": user.UUID},
		In: &infoblog.In{
			Field: "uuid",
			Args:  uuids,
		},
	}

	return e.eventRep.Updatex(ctx, condition, "shown", infoblog.Event{Shown: types.NewNullBool(true)})
}

func (e *Event) GetMy(ctx context.Context, req request.LimitOffsetReq) ([]request.NotificationResponse, error) {
	user, err := extractUser(ctx)
	if err != nil {
		return nil, err
	}
	condition := infoblog.Condition{
		Equal: &sq.Eq{"notified": true, "to_user": user.UUID},
		Order: &infoblog.Order{
			Field: "created_at",
			Asc:   false,
		},
		LimitOffset: &infoblog.LimitOffset{
			Offset: req.Offset,
			Limit:  req.Limit,
		},
	}

	events, err := e.eventRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	eventsData := make([]request.EventData, 0, len(events))
	err = e.MapStructs(&eventsData, &events)
	if err != nil {
		return nil, err
	}

	notificationRes := make([]request.NotificationResponse, len(events))
	userUUIDs := make([]interface{}, 0, len(events))
	for i := range events {
		userUUIDs = append(userUUIDs, events[i].UserUUID)
		notificationRes[i].Event = eventsData[i]
	}
	condition = infoblog.Condition{
		In: &infoblog.In{
			Field: "uuid",
			Args:  userUUIDs,
		},
	}

	users, err := e.user.userRepository.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}
	var usersData []request.UserData
	err = e.MapStructs(&usersData, &users)
	if err != nil {
		return nil, err
	}
	usersMap := make(map[string]*request.UserData, len(usersData))
	for i := range usersData {
		usersMap[usersData[i].UUID.String] = &usersData[i]
	}

	for i := range notificationRes {
		notificationRes[i].User = *usersMap[notificationRes[i].Event.UserUUID.String]
	}

	return notificationRes, nil
}

func (e *Event) WorkerPool() {
	var err error

	ticker := time.NewTicker(WorkerPoolDelay)
	condition := infoblog.Condition{
		Equal: &sq.Eq{"notified": nil},
		LimitOffset: &infoblog.LimitOffset{
			Offset: 0,
			Limit:  NumJobs,
		},
		ForUpdate: true,
	}

	for w := 0; w < NumJobs; w++ {
		go e.Worker()
	}

	var events []infoblog.Event
	for {
		select {
		case <-e.ctx.Done():
			return
		case <-ticker.C:
			events, err = e.eventRep.Listx(e.ctx, condition)
			if err != nil {
				e.logger.Error("retrieve events error", zap.Error(err))
				continue
			}
			_ = events
			for i := range events {
				e.ch <- events[i]
			}
		}
	}
}

type Response struct {
	User  request.UserData
	Event interface{}
}

func (e *Event) Worker() {
	var event infoblog.Event
	var err error
	for {
		select {
		case <-e.ctx.Done():
			return
		case event = <-e.ch:
			var res Response
			res.User, _ = e.user.Get(context.Background(), request.UserIDNickRequest{
				UUID: &event.UserUUID.String,
			})
			res.Event = event
			data := Payload{
				Action:   event.Type.Int64,
				ToUserID: event.ToUser.String,
				Data: Data{
					Res:    res,
					Action: event.Type.Int64,
				},
			}

			err = e.sendSignalMessage(data)
			if err != nil {
				continue
			}

			event.Notified = types.NewNullBool(true)
			err = e.eventRep.Update(context.Background(), event)
			if err != nil {
				continue
			}
		}
	}
}

func (e *Event) sendSignalMessage(data Payload) error {
	var body bytes.Buffer
	var err error
	err = e.Encode(&body, data)
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
	Action   int64  `json:"action"`
	ToUserID string `json:"toUserId"`
	Data     Data   `json:"data"`
}

type Data struct {
	Res    interface{} `json:"res"`
	Action int64       `json:"action"`
}

const (
	ActionSendChatMessage = iota + 1
	ActionGetChats
	ActionGetChatMessages
	ActionGetChatInfo
	ActionNotifyChatMessages
	ActionNotifySubscribe
	ActionNotifyLike
)

type WssRequest struct {
	Action  int    `json:"action"`
	Message string `json:"message"`
	UUID    string `json:"uuid"`
}
