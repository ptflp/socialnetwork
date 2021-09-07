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
	chatRep            infoblog.ChatRepository
	chatMessagesRep    infoblog.ChatMessagesRepository
	chatParticipantRep infoblog.ChatParticipantRepository
	Services           *Services
	*decoder.Decoder
}

func NewChatService(reps infoblog.Repositories, services *Services) *Chats {
	return &Chats{chatRep: reps.Chats, chatMessagesRep: reps.ChatMessages, chatParticipantRep: reps.ChatParticipant, Decoder: decoder.NewDecoder(), Services: services}
}

func (m *Chats) CreateChat(ctx context.Context) (infoblog.Chat, error) {
	chat := infoblog.Chat{
		UUID:   types.NewNullUUID(),
		Type:   types.NewNullInt64(TypeChatPrivate),
		Active: types.NewNullBool(true),
	}
	err := m.chatRep.Create(ctx, chat)

	if err != nil {
		return infoblog.Chat{}, err
	}

	return chat, nil
}

func (m *Chats) SendMessage(ctx context.Context, req request.SendMessage, chatUUID string) error {
	user, err := extractUser(ctx)
	if err != nil {
		return err
	}

	chatMessage := infoblog.ChatMessages{
		UUID:     types.NewNullUUID(),
		ChatUUID: types.NewNullUUID(chatUUID),
		UserUUID: user.UUID,
		Active:   types.NewNullBool(true),
		Message:  req.Message,
	}

	return m.chatMessagesRep.Create(ctx, chatMessage)
}

func (m *Chats) SendMessagePrivate(ctx context.Context, req request.SendMessage) error {
	var err error
	var user infoblog.User
	user, err = extractUser(ctx)
	if err != nil {
		return err
	}

	var cp []infoblog.ChatParticipant
	cp, err = m.GetPrivateParticipants(ctx, req)
	if err != nil {
		var chat infoblog.Chat
		chat, err = m.CreateChat(ctx)
		if err != nil {
			return err
		}
		err = m.AddParticipant(ctx, chat, user, infoblog.User{UUID: types.NewNullUUID(req.ToUUID)})
		if err != nil {
			return err
		}
	} else {
		if len(cp) < 2 {
			return fmt.Errorf("error count of chat participants %d", len(cp))
		}
	}

	err = m.SendMessage(ctx, req, cp[0].ChatUUID.String)
	if err != nil {
		return err
	}

	return nil
}

func (m *Chats) AddParticipant(ctx context.Context, chat infoblog.Chat, users ...infoblog.User) error {
	var err error
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

func (m *Chats) GetPrivateParticipants(ctx context.Context, req request.SendMessage) ([]infoblog.ChatParticipant, error) {
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
			Args:  []interface{}{user.UUID, types.NewNullUUID(req.ToUUID)},
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

}
