package services

import (
	"context"
	"fmt"

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

func (m *Chats) SendMessage(ctx context.Context, req request.SendMessageReq) (request.MessageData, error) {
	user, err := extractUser(ctx)
	if err != nil {
		return request.MessageData{}, err
	}

	chat := infoblog.Chat{UUID: types.NewNullUUID(req.ChatUUID)}

	chat, err = m.chatRep.Find(ctx, chat)
	if err != nil {
		return request.MessageData{}, err
	}

	chatMessage := infoblog.ChatMessages{
		UUID:     types.NewNullUUID(),
		ChatUUID: types.NewNullUUID(req.ChatUUID),
		UserUUID: user.UUID,
		Active:   types.NewNullBool(true),
		Message:  req.Message,
	}

	err = m.chatMessagesRep.Create(ctx, chatMessage)
	if err != nil {
		return request.MessageData{}, err
	}
	chat.LastMessageUUID = chatMessage.UUID
	err = m.chatRep.Update(ctx, chat)
	if err != nil {
		return request.MessageData{}, err
	}

	chatMessage, err = m.chatMessagesRep.Find(ctx, chatMessage)

	var messageData request.MessageData
	err = m.MapStructs(&messageData, &chatMessage)

	return messageData, err
}

func (m *Chats) GetMessages(ctx context.Context, req request.GetMessagesReq) ([]request.MessageData, error) {
	condition := infoblog.Condition{
		Equal: &sq.Eq{"chat_uuid": types.NewNullUUID(req.ChatUUID)},
		Order: &infoblog.Order{
			Field: "created_at",
		},
	}
	messages, err := m.chatMessagesRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}
	for i := len(messages)/2 - 1; i >= 0; i-- {
		opp := len(messages) - 1 - i
		messages[i], messages[opp] = messages[opp], messages[i]
	}
	var messagesData []request.MessageData
	err = m.MapStructs(&messagesData, &messages)
	if err != nil {
		return nil, err
	}

	return messagesData, nil
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

	if *req.UserUUID == user.UUID.String {
		return request.ChatData{}, fmt.Errorf("cant write to urself")
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
		err = m.AddParticipant(ctx, chat, user, infoblog.User{UUID: types.NewNullUUID(*req.UserUUID)})
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
	if err != nil {
		return request.ChatData{}, err
	}

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
		return request.ChatData{}, fmt.Errorf("get chat participants, no users found")
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
		Order: &infoblog.Order{
			Field: "created_at",
		},
		LimitOffset: &infoblog.LimitOffset{
			Offset: 0,
			Limit:  2,
		},
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

func (m *Chats) GetChats(ctx context.Context, req request.GetChatsReq) ([]request.ChatData, error) {
	condition := infoblog.Condition{
		Equal: &sq.Eq{"user_uuid": types.NewNullUUID(req.UserUUID)},
	}

	cpr, err := m.chatParticipantRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	chatsUUIDs := make([]interface{}, 0, len(cpr))
	for i := range cpr {
		chatsUUIDs = append(chatsUUIDs, cpr[i].ChatUUID)
	}

	condition = infoblog.Condition{
		In: &infoblog.In{
			Field: "uuid",
			Args:  chatsUUIDs,
		},
	}

	chats, err := m.chatRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	var chatsData []request.ChatData
	err = m.MapStructs(&chatsData, &chats)
	if err != nil {
		return nil, err
	}

	uuidMapsChats := make(map[string]*request.ChatData, len(chatsData))
	messagesMapsChats := make(map[string]*request.ChatData, len(chatsData))
	messagesUUIDs := make([]interface{}, 0, len(chatsData))
	for i := range chatsData {
		uuidMapsChats[chatsData[i].UUID.String] = &chatsData[i]
		messagesMapsChats[chatsData[i].LastMessageUUID.String] = &chatsData[i]
		messagesUUIDs = append(messagesUUIDs, chatsData[i].LastMessageUUID)
	}

	condition = infoblog.Condition{
		In: &infoblog.In{
			Field: "chat_uuid",
			Args:  chatsUUIDs,
		},
	}

	chatParticipants, err := m.chatParticipantRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	userMapsChats := make(map[string]*request.ChatData, len(chatParticipants))
	userUUIDs := make([]interface{}, 0, len(chatParticipants))
	for _, v := range chatParticipants {
		userUUIDs = append(userUUIDs, v.UserUUID)
		userMapsChats[v.UserUUID.String] = uuidMapsChats[v.ChatUUID.String]
	}

	if len(userUUIDs) < 1 {
		return nil, fmt.Errorf("get chat participants, no users found")
	}

	condition = infoblog.Condition{
		In: &infoblog.In{
			Field: "uuid",
			Args:  userUUIDs,
		},
	}

	users, err := m.services.User.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	usersData := make([]request.UserData, 0, len(users))

	err = m.MapStructs(&usersData, &users)
	if err != nil {
		return nil, err
	}

	for i := range usersData {
		userMapsChats[usersData[i].UUID.String].Participants = append(userMapsChats[usersData[i].UUID.String].Participants, usersData[i])
	}

	condition = infoblog.Condition{
		Order: &infoblog.Order{
			Field: "created_at",
		},
		In: &infoblog.In{
			Field: "uuid",
			Args:  messagesUUIDs,
		},
	}

	chatMessages, err := m.chatMessagesRep.Listx(ctx, condition)

	if err != nil {
		return nil, err
	}

	var messagesData []request.MessageData
	err = m.MapStructs(&messagesData, &chatMessages)

	if err != nil {
		return nil, err
	}

	for i := range messagesData {
		messagesMapsChats[messagesData[i].UUID.String].LastMessages = append(messagesMapsChats[messagesData[i].UUID.String].LastMessages, messagesData[i])
	}

	return chatsData, nil
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
			Active:     types.NewNullBool(true),
		}

		err = m.chatPrivateUsersRep.Create(ctx, cpur)
		if err != nil {
			return err
		}
		cpur = infoblog.ChatPrivateUser{
			UserUUID:   users[1].UUID,
			ToUserUUID: users[0].UUID,
			ChatUUID:   chat.UUID,
			Active:     types.NewNullBool(true),
		}
		err = m.chatPrivateUsersRep.Create(ctx, cpur)
		if err != nil {
			return err
		}
	}

	for i := range users {
		cp := infoblog.ChatParticipant{
			ChatUUID: chat.UUID,
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
