package services

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"gitlab.com/InfoBlogFriends/server/types"

	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/decoder"
	"gitlab.com/InfoBlogFriends/server/request"
)

const (
	TypeChat = iota
	TypeChatPrivate
	TypeChatGroup
)

type Chats struct {
	chatRep             infoblog.ChatRepository
	chatMessagesRep     infoblog.ChatMessagesRepository
	chatParticipantRep  infoblog.ChatParticipantRepository
	chatPrivateUsersRep infoblog.ChatPrivateUsersRepository
	services            *Services
	*decoder.Decoder
}

func NewChatService(reps infoblog.Repositories, services *Services) *Chats {
	return &Chats{chatRep: reps.Chats, chatMessagesRep: reps.ChatMessages, chatParticipantRep: reps.ChatParticipant, Decoder: decoder.NewDecoder(), services: services, chatPrivateUsersRep: reps.ChatPrivateUser}
}

func (m *Chats) CreateChat(ctx context.Context, chatType int64) (infoblog.Chat, error) {
	user, err := extractUser(ctx)
	if err != nil {
		return infoblog.Chat{}, err
	}
	chat := infoblog.Chat{
		UUID:     types.NewNullUUID(),
		Type:     types.NewNullInt64(chatType),
		Active:   types.NewNullBool(true),
		UserUUID: user.UUID,
	}
	err = m.chatRep.Create(ctx, chat)

	if err != nil {
		return infoblog.Chat{}, err
	}

	return chat, nil
}

func (m *Chats) SendMessage(ctx context.Context, req request.SendMessageReq) error {
	user, err := extractUser(ctx)
	if err != nil {
		return err
	}

	chatMessage := infoblog.ChatMessages{
		UUID:     types.NewNullUUID(),
		ChatUUID: types.NewNullUUID(req.ChatUUID),
		UserUUID: user.UUID,
		Active:   types.NewNullBool(true),
		Message:  req.Message,
	}

	return m.chatMessagesRep.Create(ctx, chatMessage)
}

func (m *Chats) Info(ctx context.Context, req request.GetInfoReq) (request.ChatData, error) {
	if req.UserUUID != nil {
		return m.GetInfoByUser(ctx, req)
	}
	if req.ChatUUID != nil {
		return m.GetInfo(ctx, req)
	}

	return request.ChatData{}, fmt.Errorf("bad request")
}

func (m *Chats) GetInfoByUser(ctx context.Context, req request.GetInfoReq) (request.ChatData, error) {
	user, err := extractUser(ctx)
	if err != nil {
		return request.ChatData{}, err
	}

	if req.UserUUID == nil {
		return request.ChatData{}, fmt.Errorf("bad request user uuid is nil")
	}
	condition := infoblog.Condition{
		Equal: &sq.Eq{"user_uuid": user.UUID, "to_user_uuid": types.NewNullUUID(*req.UserUUID)},
	}

	cpur, err := m.chatPrivateUsersRep.Listx(ctx, condition)
	if err != nil {
		return request.ChatData{}, err
	}

	var chat infoblog.Chat
	if len(cpur) < 1 {
		chat, err = m.CreateChat(ctx, TypeChatPrivate)
		if err != nil {
			return request.ChatData{}, err
		}
		err = m.AddParticipant(ctx, chat, user, infoblog.User{})
		if err != nil {
			return request.ChatData{}, err
		}
	}
	if len(cpur) > 0 {
		chat, err = m.chatRep.Find(ctx, infoblog.Chat{UUID: cpur[0].ChatUUID})
		if err != nil {
			return request.ChatData{}, err
		}
	}
	var chatData request.ChatData
	err = m.MapStructs(&chatData, &chat)

	condition = infoblog.Condition{
		Equal: &sq.Eq{"chat_uuid": chatData.UUID},
	}

	chatParticipants, err := m.chatParticipantRep.Listx(ctx, condition)
	if err != nil {
		return request.ChatData{}, err
	}
	userUUIDs := make([]interface{}, 0, 2)

	for _, v := range chatParticipants {
		userUUIDs = append(userUUIDs, v.UserUUID)
	}

	if len(userUUIDs) < 1 {
		return request.ChatData{}, fmt.Errorf("getFeedUsers no users found")
	}

	condition = infoblog.Condition{
		In: &infoblog.In{
			Field: "uuid",
			Args:  userUUIDs,
		},
	}

	users, err := m.services.User.Listx(ctx, condition)
	if err != nil {
		return request.ChatData{}, err
	}

	usersData := make([]request.UserData, 0, len(users))

	err = m.MapStructs(&usersData, &users)

	if err != nil {
		return request.ChatData{}, err
	}

	chatData.Participants = usersData

	condition = infoblog.Condition{
		Equal: &sq.Eq{"chat_uuid": chatData.UUID},
	}

	chatMessages, err := m.chatMessagesRep.Listx(ctx, condition)

	if err != nil {
		return request.ChatData{}, err
	}

	var messagesData []request.MessageData
	err = m.MapStructs(&messagesData, &chatMessages)

	if err != nil {
		return request.ChatData{}, err
	}

	chatData.LastMessages = messagesData

	return chatData, err
}

func (m *Chats) GetInfo(ctx context.Context, req request.GetInfoReq) (request.ChatData, error) {
	return request.ChatData{}, nil
}

func (m *Chats) AddParticipant(ctx context.Context, chat infoblog.Chat, users ...infoblog.User) error {
	var err error

	if chat.Type.Int64 == TypeChatPrivate {
		if len(users) != 2 {
			return fmt.Errorf("error users count in private chat %d", len(users))
		}
		cpur := infoblog.ChatPrivateUser{
			UserUUID:   users[0].UUID,
			ToUserUUID: users[1].UUID,
			ChatUUID:   chat.UUID,
			Active:     types.NullBool{},
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		}
		err = m.chatPrivateUsersRep.Create(ctx, cpur)
		if err != nil {
			return err
		}
		cpur = infoblog.ChatPrivateUser{
			UserUUID:   users[1].UUID,
			ToUserUUID: users[0].UUID,
			ChatUUID:   chat.UUID,
			Active:     types.NullBool{},
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		}
		err = m.chatPrivateUsersRep.Create(ctx, cpur)
		if err != nil {
			return err
		}
	}

	for i := range users {
		cp := infoblog.ChatParticipant{
			ChatUUID: types.NewNullUUID(),
			UserUUID: users[i].UUID,
			Type:     chat.Type,
			Active:   types.NewNullBool(true),
		}

		err = m.chatParticipantRep.Create(ctx, cp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Chats) GetPrivateParticipants(ctx context.Context, req request.SendMessageReq) ([]infoblog.ChatParticipant, error) {
	var err error
	var user infoblog.User
	var cp []infoblog.ChatParticipant
	user, err = extractUser(ctx)
	if err != nil {
		return nil, err
	}

	condition := infoblog.Condition{
		Equal: &sq.Eq{"type": TypeChatPrivate, "active": true},
		In: &infoblog.In{
			Field: "user_uuid",
			Args:  []interface{}{user.UUID, types.NewNullUUID(req.ChatUUID)},
		},
	}

	cp, err = m.chatParticipantRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	if len(cp) < 2 {
		return nil, fmt.Errorf("error count of chat participants %d", len(cp))
	}

	return cp, nil
}

func (m *Chats) GetPrivateMessages(ctx context.Context, req request.UUIDReq) ([]infoblog.ChatMessages, error) {
	condition := infoblog.Condition{
		Equal: &sq.Eq{"user_uuid": types.NewNullUUID(req.UUID), "type": TypeChatPrivate},
	}
	cp, err := m.chatParticipantRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	if len(cp) < 2 {
		return nil, fmt.Errorf("error count of chat participants %d", len(cp))
	}

	condition = infoblog.Condition{
		Equal: &sq.Eq{"chat_uuid": cp[0].ChatUUID},
	}

	ms, err := m.chatMessagesRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

func (m *Chats) GetChatByUser(ctx context.Context, req request.UUIDReq) {
	condition := infoblog.Condition{}
	m.chatParticipantRep.Listx(ctx, condition)
}
